#!/usr/bin/env bash
# Fetch Inside Airbnb datasets into seed_data/.
#
# Data source: https://insideairbnb.com/get-the-data/
# Licensed CC BY 4.0 — attribution required if you redistribute or publish
# anything derived from it.
#
# Usage:
#   scripts/fetch_seed_data.sh [--keep-gz] <city-path> <scrape-date> [file ...]
#
#   --keep-gz    leave gzipped files compressed on disk (saves ~4x space for
#                calendar/reviews; Go's compress/gzip can stream them directly)
#   city-path    e.g. united-states/ny/albany
#   scrape-date  e.g. 2026-06-16
#   file         one or more of: listings calendar reviews
#                neighbourhoods geojson summary-listings summary-reviews
#                (default: listings)
#
# Examples:
#   scripts/fetch_seed_data.sh united-states/ny/albany 2026-06-16
#   scripts/fetch_seed_data.sh united-states/ny/new-york-city 2026-06-14 listings calendar
#
# Files land in seed_data/<city>/<file>.csv, decompressed.

set -euo pipefail

BASE_URL="https://data.insideairbnb.com"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEST_ROOT="${REPO_ROOT}/seed_data"

usage() { sed -n '2,25p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'; exit 1; }

KEEP_GZ=0
if [ "${1:-}" = "--keep-gz" ]; then KEEP_GZ=1; shift; fi

[ $# -ge 2 ] || usage

CITY_PATH="$1"; shift
SCRAPE_DATE="$1"; shift
FILES=("$@")
[ ${#FILES[@]} -gt 0 ] || FILES=(listings)

CITY_SLUG="${CITY_PATH##*/}"
DEST="${DEST_ROOT}/${CITY_SLUG}"

# Map friendly name -> remote path relative to the scrape date.
remote_path_for() {
  case "$1" in
    listings)         echo "data/listings.csv.gz" ;;
    calendar)         echo "data/calendar.csv.gz" ;;
    reviews)          echo "data/reviews.csv.gz" ;;
    summary-listings) echo "visualisations/listings.csv" ;;
    summary-reviews)  echo "visualisations/reviews.csv" ;;
    neighbourhoods)   echo "visualisations/neighbourhoods.csv" ;;
    geojson)          echo "visualisations/neighbourhoods.geojson" ;;
    *) echo "unknown file type: $1" >&2; return 1 ;;
  esac
}

mkdir -p "${DEST}"

for name in "${FILES[@]}"; do
  rel="$(remote_path_for "${name}")"
  url="${BASE_URL}/${CITY_PATH}/${SCRAPE_DATE}/${rel}"

  # Preserve the real extension so we know whether to gunzip.
  case "${rel}" in
    # "geojson" is the friendly name but neighbourhoods is the real subject,
    # so don't emit geojson.geojson.
    *.geojson) out="${DEST}/neighbourhoods.geojson" ;;
    *)         out="${DEST}/${name}.csv" ;;
  esac
  if [ "${KEEP_GZ}" -eq 1 ] && [[ "${rel}" == *.gz ]]; then
    out="${out}.gz"
  fi

  echo "==> ${CITY_SLUG}/${name}"
  echo "    ${url}"

  tmp="$(mktemp)"
  # --fail so a 403/404 is an error rather than a saved error page.
  if ! curl --fail --location --silent --show-error --max-time 300 \
            --retry 3 --retry-delay 2 \
            --output "${tmp}" "${url}"; then
    echo "    FAILED — skipping" >&2
    rm -f "${tmp}"
    continue
  fi

  if [[ "${rel}" == *.gz ]] && [ "${KEEP_GZ}" -eq 0 ]; then
    gunzip --stdout "${tmp}" > "${out}"
  else
    mv "${tmp}" "${out}"
  fi
  rm -f "${tmp}"

  # Record provenance so we can tell later which scrape a file came from.
  printf '%s\t%s\t%s\n' "${name}" "${url}" "$(date -Iseconds)" \
    >> "${DEST}/SOURCE.tsv"

  # wc -l on a still-compressed file is meaningless, so stream it out first.
  if [[ "${out}" == *.gz ]]; then
    rows=$(gunzip --stdout "${out}" | wc -l)
  else
    rows=$(wc -l < "${out}")
  fi
  printf "    saved %s (%s on disk, %s lines)\n" \
    "${out#"${REPO_ROOT}/"}" "$(du -h "${out}" | cut -f1)" "${rows}"
done

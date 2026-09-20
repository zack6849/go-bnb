import { useState } from "react"
import { ListingResultView } from "@/components/listings/ListingResultView.tsx"
import type { Listing } from "@/types/Listing.ts"
import ListingSearchForm from "@/components/listings/ListingSearchForm.tsx"
import type { City } from "@/types/City.ts"

async function fetchListings(
  term: string,
  latitude: string,
  longitude: string
) {
  const params = new URLSearchParams()
  params.append("search", term)
  params.append("latitude", latitude)
  params.append("longitude", longitude)
  const response = await fetch(`/api/listings/search?${params}`)
  return await response.json()
}

export default function ListingSearch() {
  function handleSearch(term: string, city: City) {
    setActiveCity(city)
    fetchListings(
      term,
      city.Location.Latitude.toString(),
      city.Location.Longitude.toString()
    ).then((response: { results: Listing[] }) => {
      setListings(response.results)
      setShowResults(true)
    })
  }

  const [activeCity, setActiveCity] = useState<City>()
  const [showResults, setShowResults] = useState(false)
  const [listings, setListings] = useState<Listing[]>([])

  if (showResults && activeCity) {
    return (
      <ListingResultView
        listings={listings}
        city={activeCity}
        onClose={() => setShowResults(false)}
      ></ListingResultView>
    )
  } else {
    return <ListingSearchForm onSearch={handleSearch} />
  }
}

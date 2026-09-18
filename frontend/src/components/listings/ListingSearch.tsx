import {useState} from "react";
import {ListingResultView} from "@/components/listings/ListingResultView.tsx";
import type {Listing} from "@/types/Listing.ts";
import ListingSearchForm from "@/components/listings/ListingSearchForm.tsx";

async function fetchListings(term: string) {
    // const value = await geocode(term);
    // const coords = value.features[0].geometry.coordinates;
    const params = new URLSearchParams();
    params.append("search", term);
    // params.append("latitude", coords[1]);
    // params.append("longitude", coords[0]);
    params.append("latitude", "42.6511674")
    params.append("longitude", "-73.754968")
    const response = await fetch(`/api/listings?${params}`);
    return await response.json();
}

export default function ListingSearch() {

    function handleSearch(term: string) {
        fetchListings(term).then((response: { listings: Listing[] }) => {
            setListings(response.listings)
            setShowResults(true)
        })
    }

    const [showResults, setShowResults] = useState(false);
    const [listings, setListings] = useState<Listing[]>([]);

    if (showResults) {
        return <ListingResultView listings={listings}></ListingResultView>
    } else {
        return <ListingSearchForm onSearch={handleSearch}/>
    }
}
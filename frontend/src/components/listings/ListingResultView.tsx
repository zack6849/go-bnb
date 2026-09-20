import type {Listing} from "@/types/Listing.ts";
import ListingCard from "@/components/listings/ListingCard.tsx";

interface ListingResultProps {
    listings: Listing[],
}

export function ListingResultView(props: ListingResultProps) {
    let listingsComponents = props.listings.map((l => {
        return (<ListingCard key={l.ID} listing={l}/>)
    }));

    return (
        <div className={"grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3"}>
            {listingsComponents}
        </div>
    )
}
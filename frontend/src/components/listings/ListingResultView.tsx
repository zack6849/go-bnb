import type {Listing} from "@/types/Listing.ts";
import ListingCard from "@/components/listings/ListingCard.tsx";
import type {City} from "@/types/City.ts";
import {Button} from "@/components/ui/button.tsx";

interface ListingResultProps {
    listings: Listing[],
    city: City,
    onClose: () => void
}

export function ListingResultView(props: Readonly<ListingResultProps>) {
    function getHeader(location: string) {
        return (
            <div className={"flex justify-between mb-4"}>
                <h1 className={'text-3xl'}>Let's go to <b className={"text-primary"}>{location}</b> together</h1>
                <Button onClick={props.onClose}>Back</Button>
            </div>
        )
    }

    let listingsComponents = props.listings.map((l => {
        return (<ListingCard key={l.ID} listing={l}/>)
    }));

    let header = getHeader(props.city.Name)

    return (
        <div>
            {header}
            <div className={"grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3"}>
                {listingsComponents}
            </div>
        </div>
    )
}
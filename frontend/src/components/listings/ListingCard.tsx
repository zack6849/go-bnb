import {Card, CardAction, CardContent, CardFooter, CardHeader} from "@/components/ui/card.tsx";
import {Button} from "@/components/ui/button.tsx";
import HostWidget from "@/components/host/HostWidget.tsx"
import type {Listing} from "@/types/Listing.ts";

export interface ListingCardProps {
    listing: Listing,
}

export default function ListingCard(props: Readonly<ListingCardProps>) {
    const hosts = props.listing.Hosts.map((h => {
        return (<HostWidget key={h.ID} host={h}/>)
    }))
    return (
        <Card className={"flex"}>
            <CardHeader>
                <h3 className={"text-2xl"}>{props.listing.Tagline}</h3>
            </CardHeader>
            <CardContent className={"flex flex-col grow"}>
                <img src={props.listing.PictureURL} loading={"lazy"} alt={props.listing.Tagline}></img>
                <div className={"py-4 grow"}>
                    <p>{props.listing.Description}</p>
                </div>
                <h3 className={"text-2xl"}>Hosted By</h3>
                {hosts}
            </CardContent>
            <CardFooter className={"justify-between"}>
                <CardAction>
                    <Button>
                        <a target={"_blank"} href={props.listing.ListingUrl}>View Listing</a>
                    </Button>
                </CardAction>
            </CardFooter>
        </Card>
    )
}
import {Card, CardAction, CardContent, CardFooter, CardHeader} from "@/components/ui/card.tsx";
import {Button} from "@/components/ui/button.tsx";


export interface Listing {
    AirBNBID: string,
    ID: string,
    Tagline: string,
    Description: string,
    ListingUrl: string,
    NumBeds: number,
    NumBaths: number,
}

export interface ListingCardProps {
    listing: Listing,
    className?: string
}

export default function ListingCard(props: Readonly<ListingCardProps>) {
    return (
        <Card>
            <CardHeader>
                <h3 className={"text-2xl"}>{props.listing.Tagline}</h3>
            </CardHeader>
            <CardContent>
                <div className={props.className}>
                    <p>{props.listing.Description}</p>
                </div>
            </CardContent>
            <CardFooter>
                <CardAction>
                    <Button>
                        <a target={"_blank"} href={props.listing.ListingUrl}>View Listing</a>
                    </Button>
                </CardAction>
            </CardFooter>
        </Card>
    )
}
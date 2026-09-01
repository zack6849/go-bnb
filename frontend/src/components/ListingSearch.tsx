import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Field, FieldLabel} from "@/components/ui/field.tsx";
import {useState} from "react";
import ListingCard, {type Listing} from "@/components/ListingCard.tsx";
import {Card, CardContent, CardHeader} from "@/components/ui/card.tsx";


function getPlaceholderText() {
    return "Where do you want to go?"
}

function geocode(term: string) {
    const params = new URLSearchParams();
    params.append("apiKey", "f297699f0c4143af808afb6bc8cfc291")
    params.append("text", term);
    return fetch(`https://api.geoapify.com/v1/geocode/search?${params}`, {
        method: 'GET',
    }).then(response => response.json())
}

async function fetchListings(term: string) {
    const value = await geocode(term);
    const coords = value.features[0].geometry.coordinates;
    const params = new URLSearchParams();
    params.append("search", term);
    //TODO: these may be swapped, doesn't matter right now though
    params.append("latitude", coords[1]);
    params.append("longitude", coords[0]);
    const response = await fetch(`/api/listings?${params}`);
    return await response.json();
}


export default function ListingSearch() {

    const [searchTerm, setSearchTerm] = useState('');
    const [showResults, setShowResults] = useState(false);
    const [listings, setListings] = useState<Listing[]>([]);

    const handleButton = () => {
        fetchListings(searchTerm).then((response: { listings: Listing[] }) => {
            setListings(response.listings)
            setShowResults(true)
        })
    };


    if (showResults) {
        const listingsElements = listings?.map((l => {
            return (<ListingCard className={"mx-2"} key={l.ID} listing={l}/>)
        }))
        return (
            <Card>
                <CardHeader>
                    <h1 className={"text-3xl"}>
                        Let's go to <b className={'text-primary'}>{searchTerm}</b> together
                    </h1>
                </CardHeader>
                <CardContent>
                    {listingsElements}
                </CardContent>
            </Card>
        )
    }

    return (
        <Card>
            <CardHeader>
                <h1 className={"text-3xl"}>
                    Let's go <b className={'text-primary'}>somewhere</b> together
                </h1>
            </CardHeader>
            <CardContent>
                <div className={"p-4"}>
                    <div></div>
                    <Field>
                        <FieldLabel htmlFor="input-button-group">Search</FieldLabel>
                    </Field>
                    <Field className={"my-2"}>
                        <Input
                            placeholder={getPlaceholderText()}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            value={searchTerm}
                        />
                    </Field>
                    <div>
                        <Button onClick={handleButton}>Search</Button>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
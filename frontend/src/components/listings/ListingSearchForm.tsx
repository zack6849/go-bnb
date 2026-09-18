import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Field, FieldLabel} from "@/components/ui/field.tsx";
import {useState} from "react";
import {Card, CardContent, CardHeader} from "@/components/ui/card.tsx";
import * as React from "react";


function getHeader(location: string) {
    return (
        <CardHeader>
            <h1 className={'text-3xl'}>Let's go <b className={"text-primary"}>{location}</b> together</h1>
        </CardHeader>
    )
}

function getPlaceholderText() {
    return "Where do you want to go?"
}

interface ListingSearchFormProps {
    onSearch: (term: string) => void;
}

export default function ListingSearchForm({onSearch}: Readonly<ListingSearchFormProps>) {

    const [searchTerm, setSearchTerm] = useState('');
    const location = searchTerm ? `to ${searchTerm}` : `somewhere`;
    const header = getHeader(location)

    let content: React.JSX.Element[] | React.JSX.Element = (
        <div className={"p-4"}>
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
                <Button onClick={() => onSearch(searchTerm)}>Let's go places</Button>
            </div>
        </div>
    )

    return (
        <Card>
            {header}
            <CardContent>
                {content}
            </CardContent>
        </Card>
    )
}
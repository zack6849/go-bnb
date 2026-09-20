import {Card, CardContent, CardHeader} from "@/components/ui/card.tsx"
import * as React from "react"
import {useAsyncList} from 'react-aria-components/useAsyncList'
import type {City} from "@/types/City.ts"
import {
    Combobox,
    ComboboxContent,
    ComboboxEmpty,
    ComboboxInput,
    ComboboxItem,
    ComboboxList,
} from "@/components/ui/combobox"

interface ListingSearchFormProps {
    onSearch: (term: string, city: City) => void
}

function getPlaceholderText() {
    return "Where do you want to go?"
}

export default function ListingSearchForm(props: Readonly<ListingSearchFormProps>) {
    function getHeader(location: string) {
        return (
            <h1 className={'text-3xl'}>Let's go <b className={"text-primary"}>{location}</b> together</h1>
        )
    }

    let list = useAsyncList<City>({
        async load({signal, filterText}) {
            let params = new URLSearchParams()
            params.set("search", filterText ?? "")
            let res = await fetch(`/api/cities/search?${params}`, {signal})
            let json = await res.json()
            return {
                items: json.results
            }
        }
    })


    function handleSearchUpdate(newValue: string, city: City) {
        props.onSearch(newValue, city)
    }


    const location = list.filterText ? `to ${list.filterText}` : `somewhere`
    const items = list.items.map((c: City) => {
        return (
            <ComboboxItem key={c.ID} onClick={() => handleSearchUpdate(c.Name, c)}>
                {c.Name}
            </ComboboxItem>
        )
    })

    const header = getHeader(location)
    let content: React.JSX.Element[] | React.JSX.Element = (
        <div className={"p-4"}>
            <Combobox<City>
                items={list.items}
                inputValue={list.filterText}
                onInputValueChange={list.setFilterText}
            >
                <ComboboxInput placeholder={getPlaceholderText()} value={list.filterText}/>
                <ComboboxEmpty>
                    No cities
                </ComboboxEmpty>
                <ComboboxContent>
                    <ComboboxList>
                        {items}
                    </ComboboxList>
                </ComboboxContent>
            </Combobox>
        </div>
    )

    return (
        <Card>
            <CardHeader>
                {header}
            </CardHeader>
            <CardContent>
                {content}
            </CardContent>
        </Card>
    )
}
import type {Host} from "@/types/Host.ts";

export interface Listing {
    AirBNBID: string,
    ID: string,
    Tagline: string,
    Description: string,
    ListingUrl: string,
    NumBeds: number,
    NumBaths: number,
    Hosts: Host[],
    PictureURL: string,
}
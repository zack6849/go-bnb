import type { Host } from "@/types/Host.ts"

export interface Listing {
  AirBnbID: string
  ID: string
  Tagline: string
  Description: string
  ListingURL: string
  NumBeds: number
  NumBaths: number
  Hosts: Host[]
  PictureURL: string
}

export interface City {
    ID: string
    Name: string,
    Location: {
        Latitude: number,
        Longitude: number,
    }
    Subdivision: string
    FullName: string,
    Timezone: string,
    Hash: string,
}
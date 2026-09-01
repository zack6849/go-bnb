
import {HeaderConfig} from "@/config/header.ts";
import Navbar from "@/components/ui/navbar"
import Footer from "@/components/ui/footer.tsx";
import ListingSearch from "@/components/ListingSearch.tsx";

export function App() {
    return (
        <div className={"bg-background text-foreground flex flex-col flex-1"}>
            <Navbar navItems={HeaderConfig}></Navbar>
            <main className={"grow p-10"}>
                <ListingSearch/>
            </main>
            <Footer/>
        </div>
    )
}

export default App

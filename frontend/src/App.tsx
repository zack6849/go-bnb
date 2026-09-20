import { HeaderConfig } from "@/config/header.ts"
import Navbar from "@/components/ui/navbar"
import Footer from "@/components/ui/footer.tsx"
import ListingSearch from "@/components/listings/ListingSearch.tsx"

export function App() {
  return (
    <div className={"flex flex-1 flex-col bg-background text-foreground"}>
      <Navbar navItems={HeaderConfig}></Navbar>
      <main className={"grow p-10"}>
        <ListingSearch />
      </main>
      <Footer />
    </div>
  )
}

export default App

import {
    NavigationMenu,
    NavigationMenuItem, NavigationMenuLink,
    NavigationMenuList,
} from "@/components/ui/navigation-menu.tsx";
import {Card, CardContent} from "@/components/ui/card.tsx";

export interface NavItem {
    label: string;
    link?: string;
    icon?: string;
    items?: NavItem[]
}

interface NavProps {
    navItems: NavItem[]
}

export default function Navbar({navItems}: Readonly<NavProps>) {
    const navItem = navItems.map(navItem => {
        return (
            <NavigationMenuItem key={navItem.label}>
                <NavigationMenuLink href={navItem.link}>
                    {navItem.label}
                </NavigationMenuLink>
            </NavigationMenuItem>
        )
    });
    return (
        <nav>
            <Card className={"background-sidebar-background"}>
                <CardContent>
                    <NavigationMenu className={"flex grow-1"}>
                        <span className={""}>GoBNB</span>
                        <NavigationMenuList className={"gap-8 justify-start"}>
                            {navItem}
                        </NavigationMenuList>
                    </NavigationMenu>
                </CardContent>
            </Card>

        </nav>
    )
}
import { Component } from "react"
import { Card, CardContent } from "@/components/ui/card.tsx"

export default class Footer extends Component<any, any> {
  render() {
    return (
      <footer>
        <Card>
          <CardContent>&copy; GoBNB 2026</CardContent>
        </Card>
      </footer>
    )
  }
}

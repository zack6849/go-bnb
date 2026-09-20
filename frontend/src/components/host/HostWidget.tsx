import {
  Avatar,
  AvatarBadge,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar.tsx"
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card.tsx"
import { LucideMapPin, LucideStar } from "lucide-react"
import type { Host } from "@/types/Host.ts"

export interface HostProp {
  host: Host
}

function HostAvatar(props: Readonly<HostProp>) {
  const isSuperHost = props.host.SuperHost
  let badge = null
  if (isSuperHost) {
    badge = (
      <AvatarBadge>
        <LucideStar />
      </AvatarBadge>
    )
  }
  return (
    <Avatar>
      <AvatarImage
        src={props.host.PictureURL}
        alt={props.host.Name + "'s profile picture"}
      />
      ,<AvatarFallback>{props.host.Name}</AvatarFallback>
      {badge}
    </Avatar>
  )
}

export default function HostWidget(props: Readonly<HostProp>) {
  return (
    <HoverCard>
      <HoverCardTrigger>
        <a href={props.host.Slug}>
          <div className={"flex w-fit flex-col items-center"}>
            <HostAvatar host={props.host}></HostAvatar>
            <div className={"mt-2"}>{props.host.Name}</div>
          </div>
        </a>
      </HoverCardTrigger>
      <HoverCardContent side="bottom" align="start">
        <div>
          <HostAvatar host={props.host}></HostAvatar>
          <h1>{props.host.Name}</h1>
        </div>
        <div>
          <LucideMapPin /> {props.host.Location}
        </div>
        <div>
          <span>{props.host.Description}</span>
        </div>
      </HoverCardContent>
    </HoverCard>
  )
}

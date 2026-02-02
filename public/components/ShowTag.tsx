import "./ShowTag.scss"

import React from "react"
import { Tag } from "@fider/models"
import { classSet } from "@fider/services"
import EyeSlash from "@fider/assets/images/heroicons-eyeslash.svg"
import TagSolid from "@fider/assets/images/heroicons-tagsolid.svg"
import { Icon } from "./common"

interface TagProps {
  tag: Tag
  circular?: boolean
  link?: boolean
}

// const textColor = (color: string) => {
//   const components = getRGB(color)
//   const bgDelta = components.R * 0.299 + components.G * 0.587 + components.B * 0.114
//   return bgDelta > 140 ? "#333" : "#fff"
// }

export const ShowTag = (props: TagProps) => {
  const className = classSet({
    "c-tag": true,
    "c-tag--circular": props.circular === true,
  })

  const title = `${props.tag.name}${props.tag.isPublic ? "" : " (Private)"}`
  const content = (
    <>
      <Icon style={{ color: `#${props.tag.color}` }} className="pr-1" height="18" width="18" sprite={TagSolid}></Icon>
      {!props.tag.isPublic && !props.circular && <Icon height="14" width="14" sprite={EyeSlash} className="mr-1" />}
      {props.circular ? "" : props.tag.name || "Tag"}
    </>
  )

  // Use span when not a link to avoid invalid nested <a> (e.g. inside post list link)
  if (!props.link || !props.tag.slug) {
    return (
      <span title={title} className={className}>
        {content}
      </span>
    )
  }

  return (
    <a href={`/?tags=${props.tag.slug}`} title={title} className={className}>
      {content}
    </a>
  )
}

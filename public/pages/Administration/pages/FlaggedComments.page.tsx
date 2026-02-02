import React, { useEffect, useState } from "react"
import { AdminPageContainer } from "@fider/pages/Administration/components/AdminBasePage"
import { Button, Avatar, UserName } from "@fider/components"
import { http } from "@fider/services"
import { Comment } from "@fider/models"
import { HStack, VStack } from "@fider/components/layout"
import { Trans, useLingui } from "@lingui/react/macro"
import { formatDate } from "@fider/services"

interface FlaggedCommentItem {
  comment: Comment
  postNumber: number
  postTitle: string
  postSlug: string
  flagsCount: number
}

const FlaggedCommentsPage = () => {
  const { i18n } = useLingui()
  const [items, setItems] = useState<FlaggedCommentItem[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchFlagged = async () => {
      const result = await http.get<FlaggedCommentItem[]>("/api/v1/admin/comments/flagged")
      if (result.ok) {
        setItems(result.data)
      }
      setLoading(false)
    }
    fetchFlagged()
  }, [])

  return (
    <AdminPageContainer
      id="p-admin-flagged"
      name="flagged"
      title="Flagged Comments"
      subtitle="Comments flagged by users for inappropriateness"
    >
      {loading ? (
            <p className="text-muted">
              <Trans id="showpost.loading">Loading...</Trans>
            </p>
          ) : items.length === 0 ? (
            <p className="text-muted">
              <Trans id="admin.flagged.empty">No flagged comments.</Trans>
            </p>
          ) : (
            <VStack spacing={4}>
              {items.map((item) => (
                <div key={item.comment.id} className="p-4 border rounded-md bg-white shadow-sm">
                  <HStack justify="between" align="start">
                    <VStack spacing={2} align="start">
                      <HStack spacing={2} align="center">
                        <span className="text-xs px-2 py-0.5 rounded bg-yellow-100 text-yellow-800">
                          {(i18n as any)._({ id: "label.flagcount", message: "{count} flag(s)" }, { count: item.flagsCount })}
                        </span>
                        <a href={`/posts/${item.postNumber}/${item.postSlug}#comment-${item.comment.id}`} className="text-sm font-medium text-green-700 hover:underline">
                          {item.postTitle} #{item.postNumber}
                        </a>
                      </HStack>
                      <p className="text-sm text-gray-700">{item.comment.content.slice(0, 200)}{item.comment.content.length > 200 ? "…" : ""}</p>
                      <HStack spacing={2} align="center">
                        <Avatar user={item.comment.user} size="small" />
                        <UserName user={item.comment.user} />
                        <span className="text-xs text-gray-400">{formatDate("en", item.comment.createdAt)}</span>
                      </HStack>
                    </VStack>
                    <Button variant="secondary" size="small" href={`/posts/${item.postNumber}/${item.postSlug}#comment-${item.comment.id}`}>
                      <Trans id="action.view">View</Trans>
                    </Button>
                  </HStack>
                </div>
              ))}
            </VStack>
          )}
    </AdminPageContainer>
  )
}

export default FlaggedCommentsPage

import React, { useEffect, useState } from "react"
import { AdminPageContainer } from "@fider/pages/Administration/components/AdminBasePage"
import { Button, Avatar, UserName } from "@fider/components"
import { http } from "@fider/services"
import { Post } from "@fider/models"
import { HStack, VStack } from "@fider/components/layout"
import { Trans, useLingui } from "@lingui/react/macro"
import { formatDate } from "@fider/services"

interface FlaggedPostItem {
  post: Post
  flagsCount: number
}

const FlaggedPostsPage = () => {
  const { i18n } = useLingui()
  const [items, setItems] = useState<FlaggedPostItem[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchFlagged = async () => {
      const result = await http.get<FlaggedPostItem[]>("/api/v1/admin/posts/flagged")
      if (result.ok && Array.isArray(result.data)) {
        setItems(result.data.filter((item) => item && item.post && item.post.id != null))
      }
      setLoading(false)
    }
    fetchFlagged()
  }, [])

  return (
    <AdminPageContainer
      id="p-admin-flagged-posts"
      name="flaggedPosts"
      title="Flagged Ideas"
      subtitle="Ideas flagged by users for inappropriateness"
    >
      {loading ? (
            <p className="text-muted">
              <Trans id="showpost.loading">Loading...</Trans>
            </p>
          ) : items.length === 0 ? (
            <p className="text-muted">
              <Trans id="admin.flaggedposts.empty">No flagged ideas.</Trans>
            </p>
          ) : (
            <VStack spacing={4}>
              {items.map((item) => {
                const post = item.post
                if (!post) return null
                return (
                  <div key={post.id} className="p-4 border rounded-md bg-white shadow-sm">
                    <HStack justify="between" align="start">
                      <VStack spacing={2} align="start">
                        <HStack spacing={2} align="center">
                          <span className="text-xs px-2 py-0.5 rounded bg-yellow-100 text-yellow-800">
                            {(i18n as any)._({ id: "label.flagcount", message: "{count} flag(s)" }, { count: item.flagsCount })}
                          </span>
                          <a href={`/posts/${post.number}/${post.slug}`} className="text-sm font-medium text-green-700 hover:underline">
                            {post.title} #{post.number}
                          </a>
                        </HStack>
                        {post.description && (
                          <p className="text-sm text-gray-700">{post.description.slice(0, 200)}{post.description.length > 200 ? "…" : ""}</p>
                        )}
                        <HStack spacing={2} align="center">
                          {post.user && (
                            <>
                              <Avatar user={post.user} size="small" />
                              <UserName user={post.user} />
                            </>
                          )}
                          <span className="text-xs text-gray-400">{formatDate("en", post.createdAt)}</span>
                        </HStack>
                      </VStack>
                      <Button variant="secondary" size="small" href={`/posts/${post.number}/${post.slug}`}>
                        <Trans id="action.view">View</Trans>
                      </Button>
                    </HStack>
                  </div>
                )
              })}
            </VStack>
          )}
    </AdminPageContainer>
  )
}

export default FlaggedPostsPage

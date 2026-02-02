import React from "react"
import { ErrorPageWrapper } from "./components/ErrorPageWrapper"
import { Trans } from "@lingui/react/macro"

interface Error500Props {
  errorMessage?: string
}

const Error500 = (props: Error500Props) => {
  return (
    <ErrorPageWrapper id="p-error500" showHomeLink={true}>
      <h1 className="text-display uppercase">
        <Trans id="error.internalerror.title">Shoot! Well, this is unexpected…</Trans>
      </h1>
      <p>
        <Trans id="error.internalerror.text">An error has occurred and we&apos;re working to fix the problem! We’ll be up and running shortly.</Trans>
      </p>
      {props.errorMessage && (
        <pre className="mt-4 p-4 bg-gray-100 text-left text-sm overflow-auto rounded border border-gray-300">
          {props.errorMessage}
        </pre>
      )}
    </ErrorPageWrapper>
  )
}

export default Error500

import React from "react"
import { SignInControl, TenantLogo, LegalNotice } from "@fider/components"
import "./SignIn.page.scss"
import { Trans } from "@lingui/react/macro"
import { useFider } from "@fider/hooks"

const Locked = (): JSX.Element => {
  const fider = useFider()
  return (
    <>
      <p className="text-title">
        <Trans id="signin.message.locked.title">
          <strong>{fider.session.tenant.name}</strong> is currently locked.
        </Trans>
      </p>
      <Trans id="signin.message.locked.text">To reactivate this site, sign in with an administrator account and update the required settings.</Trans>
    </>
  )
}

const Private = (): JSX.Element => {
  const fider = useFider()
  return (
    <>
      <p className="text-title">
        <Trans id="signin.message.private.title">
          <strong>{fider.session.tenant.name}</strong> is a private space, you must sign in to participate and vote.
        </Trans>
      </p>
      <Trans id="signin.message.private.text">If you have an account or an invitation, you may use following options to sign in.</Trans>
    </>
  )
}

export const SignInPage = () => {
  const fider = useFider()

  const onCodeVerified = () => {
    // User is authenticated - redirect to the appropriate URL
    const redirect = new URLSearchParams(window.location.search).get("redirect")
    if (redirect && redirect.startsWith("/")) {
      location.href = fider.settings.baseURL + redirect
    } else {
      location.href = fider.settings.baseURL
    }
  }

  const getRedirectToUrl = () => {
    const redirect = new URLSearchParams(window.location.search).get("redirect")

    if (redirect && redirect.startsWith("/")) {
      return fider.settings.baseURL + redirect
    }

    return fider.settings.baseURL
  }

  // Corrected: Access casEnabled directly from fider.settings
  const casEnabled = fider.settings.casEnabled === true;
  const currentRedirect = new URLSearchParams(window.location.search).get("redirect") || "/";
  const casLoginUrl = `/cas/login?redirect=${encodeURIComponent(currentRedirect)}`;


  return (
    <div id="p-signin" className="page container w-max-6xl">
      <div className="h-20 text-center mb-4">
        <TenantLogo size={100} />
      </div>
      <div className="text-center w-max-4xl mx-auto mb-4">{fider.session.tenant.isPrivate ? <Private /> : <Locked />}</div>

      {casEnabled && (
        <div className="flex justify-center mb-4">
          <a
            href={casLoginUrl}
            className="btn btn-block btn-social bg-uab-green hover:bg-uab-green-dark text-white"
            style={{ backgroundColor: '#1E6B52', borderColor: '#1E6B52' }}
          >
            <Trans id="signin.uab.cas">Sign in with UAB</Trans>
          </a>
        </div>
      )}

      <SignInControl onCodeVerified={onCodeVerified} useEmail={true} redirectTo={getRedirectToUrl()} />
      <LegalNotice />
    </div>
  )
}

export default SignInPage

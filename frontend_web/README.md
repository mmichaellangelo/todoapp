# Todo Frontend Web App
## About
This is the frontend web interface for my todo app. Powered by [SvelteKit](kit.svelte.dev), which I am currently learning how to use.
## Routes

- **/** << ROOT - landing page or redirect to dashboard if logged in >>
- **/login** << LOGIN >>
- **/dashboard** << DASHBOARD HOME >>
    - **/dashboard/overview** << OVERVIEW >>

## State management
Username and UserID are parsed from access token when user logs in and on every subsequent request via server hook. 
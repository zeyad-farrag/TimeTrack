export default function Home() {
  return (
    <div className="flex flex-1 bg-background font-sans text-foreground">
      <main className="mx-auto flex min-h-screen w-full max-w-5xl flex-col justify-center gap-10 px-6 py-16 sm:px-10">
        <div className="flex max-w-2xl flex-col gap-5">
          <p className="text-sm font-medium uppercase text-muted-foreground">
            TimeTrack
          </p>
          <h1 className="text-4xl font-semibold leading-tight tracking-normal sm:text-5xl">
            Team capacity tracking starts here.
          </h1>
          <p className="max-w-xl text-lg leading-8 text-muted-foreground">
            This scaffold is ready for the auth bridge, team views, and capacity
            workflows that land in the next stories.
          </p>
        </div>
        <div className="flex flex-col gap-3 text-base font-medium sm:flex-row">
          <a
            className="flex h-11 items-center justify-center rounded-md bg-primary px-4 text-primary-foreground transition-colors hover:opacity-90"
            href="/api/v1"
          >
            API placeholder
          </a>
          <a
            className="flex h-11 items-center justify-center rounded-md border border-border px-4 transition-colors hover:bg-muted"
            href="/healthz"
          >
            Health check
          </a>
        </div>
      </main>
    </div>
  );
}

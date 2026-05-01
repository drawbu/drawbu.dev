interface ImportMetaEnv {
  readonly GIT_REV?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

declare namespace astroHTML.JSX {
  interface HTMLAttributes {
    // htmx-preload
    preload?: "mousedown" | "mouseover" | "always";
  }
}

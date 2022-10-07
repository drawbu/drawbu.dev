<script>
  import Website from './lib/website.svelte'

  const websites_list = [
    {
      url: 'https://ascii.drawbu.dev',
      name: 'ASCII ART Convertor',
      description: 'Converts an image given by a user, to a new colored ASCII ART image.',
      image: 'https://raw.githubusercontent.com/drawbu/ASCII-ART-Generator/master/assets/screenshot.png',
    },
  ]
</script>

<div id="app">
  <header>
    <h1><code>{"<"}<span>drawbu.dev</span> \></code></h1>
    <code>
      I'm <a href="https://drawbu.me">drawbu</a> and welcome to my projects portfolio.
      <br/> Fell free to click on anything :)
    </code>
  </header>
  <main>
    {#each websites_list as website}
      <Website website={website} />
    {/each}
  </main>
</div>

<style lang="scss">
  #app {
    width: 800px;
    margin: 0 auto;

    @media (max-width: 840px) {
      width: auto;
    }

    header {
      margin-bottom: 3rem;
      text-align: center;
      padding: 1em;

      h1 {
        color: darkgray;
        span {
          color: dodgerblue;
        }
        font-size: 30px;
        font-weight: 100;

        @media (min-width: 350px) {
          font-size: 50px;
        }
      }

      a {
        color: dodgerblue;
      }
    }

    main {
      margin-bottom: 10px;
    }
  }

</style>
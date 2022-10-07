<script>
  import ProjectCard from './lib/project_card.svelte'

  const projects_list = [
    {
      url: 'https://ascii.drawbu.dev',
      name: 'ASCII ART Convertor',
      description: 'Converts an image given by a user, to a new colored ASCII ART image.',
      image: 'https://raw.githubusercontent.com/drawbu/ASCII-ART-Generator/master/assets/screenshot.png',
      repo: 'https://github.com/drawbu/ASCII-ART-Generator'
    },
  ]
</script>

<div id="app">
  <header>
    <h1><code>{'<'}<span>drawbu.dev</span> \></code></h1>
    <code>
      I'm <a href="https://drawbu.me">drawbu</a> and welcome to my projects portfolio.
      <br/> Fell free to click on anything :)
    </code>
  </header>
  <main>
    {#each projects_list as project}
      <ProjectCard project={project} />
    {/each}
  </main>
</div>

<style lang="scss">
  #app {
    width: 800px;
    margin: auto;

    @media (max-width: 840px) {
      width: auto;
    }

    header {
      margin-bottom: 3rem;
      text-align: center;
      padding: 8px;

      h1 {
        color: darkgray;
        font-size: 30px;
        font-weight: 100;
        margin-top: 16px;
        margin-bottom: 16px;
        transition-duration: 500ms;

        span {
          color: dodgerblue;
        }

        @media (min-width: 450px) {
          font-size: 50px;
        }
      }
    }

    main {
      margin-bottom: 10px;
    }
  }

</style>
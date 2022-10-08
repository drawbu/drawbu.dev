<script>
  import Header from './lib/header.svelte'
  import ProjectCard from './lib/project_card.svelte';
  import Footer from './lib/footer.svelte'

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
  <Header />
  <main>
    {#each projects_list as project}
      <ProjectCard project={project} />
    {/each}
  </main>
  <Footer />
</div>

<style lang="scss">
  #app {
    width: 800px;
    margin: auto;

    @media (max-width: 840px) {
      width: auto;
    }

    main {
      margin-bottom: 10px;
    }
  }
</style>

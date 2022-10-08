<script>
  import { GithubIcon } from 'svelte-feather-icons';

  export let project;
</script>

<div class="website">
  <div class="website-title">
    {#if project.repo}
      <a href={project.repo} class="repo-link">
        <GithubIcon />
      </a>
    {/if}
    <h2>{project.name}</h2>
    <a href={project.url}><code>{project.url}</code></a>
  </div>
  <p>{project.description}</p>
  <img src={project.image} alt="" />
</div>

<style lang="scss">
  .website {
    border: darkgray solid 2px;
    border-radius: 10px;
    text-align: center;
    width: 90%;
    margin: auto;
    padding: 8px;

    .website-title {
      margin-bottom: 16px;

      .repo-link {
        color: #333;
        :global(svg.feather-github) {
          float: right;
        }

        &:hover {
          color: dodgerblue;
        }
      }
      h2 {
        margin-bottom: 4px;
      }
    }

    img {
      max-height: 300px;

      @media (max-width: 840px) {
        height: auto;
        max-width: 100%;
      }
    }
  }
</style>

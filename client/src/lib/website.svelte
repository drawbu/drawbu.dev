<script>
    export let website;
</script>


<div class="website">
    <h2>{website.name}</h2>
    <a href={website.url}><code>{website.url}</code></a>
    <p>{website.description}</p>
    <img src={website.image} alt="" />
</div>


<style lang="scss">
    .website {
      border: darkgray solid 2px;
      border-radius: 10px;
      background-color: transparent;
      width: 100%;
      text-align: center;
      color: black;

      h2 {
        margin-bottom: 4px;
      }

      a {
        color: dodgerblue;
      }

      img {
        height: 300px;

        @media (max-width: 840px) {
          width: 100%;
          height: auto;
        }
      }
    }
</style>
<script>
  export let website;
</script>


<div class="website">
  <div class="website-title">
    <h2>{website.name}</h2>
    <a href={website.url}><code>{website.url}</code></a>
  </div>
  <p>{website.description}</p>
  <img src={website.image} alt="" />
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
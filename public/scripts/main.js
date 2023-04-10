const main = document.querySelector('main');
const scrollDownButton = document.getElementById('scroll-down-btn');

scrollDownButton.addEventListener('click', () => {
  main.scrollIntoView({ behavior: 'smooth' });
});

const portfolio = document.getElementById('portfolio');
const scrollDownButton = document.getElementById('scroll-down-btn');

scrollDownButton.addEventListener('click', () => {
  portfolio.scrollIntoView({ behavior: 'smooth' });
});

const header = document.querySelector('header');
const main = document.querySelector('main');
const observer = new IntersectionObserver((entries) => {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      header.classList.add('scrolled')
      observer.unobserve(main);
    }
  })
});
observer.observe(main);

const scrollDownButton = document.getElementById('scroll-down-btn');
scrollDownButton.addEventListener('click', () => {
  header.classList.add('scrolled');
});

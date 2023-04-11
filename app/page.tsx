'use client';

import Image from 'next/image'
import React from 'react';

interface CardProps {
  bg_src: string,
  name: string,
  link?: string,
  repo?: string,
}

const Card: React.FC<CardProps> = ({ bg_src, name, link, repo }) => {
  function returnRepo(): JSX.Element | undefined {
    if (!repo) {
      return;
    }
    return (<a href={repo}>
      <Image
        src="/images/github.svg"
        alt="GitHub logo"
        width={24}
        height={24}
      />
    </a>);
  }

  function returnLink(): JSX.Element | undefined {
    if (!link) {
      return;
    }
    return (<p><a href={link}>{link}</a></p>);
  }

  return (
    <div className="card">
      <Image
        src={bg_src}
        alt="Background image"
        width={300}
        height={300}
      />
      <div className="panel">
        <div className="content">
          <h3>{name}</h3>
          {returnLink()}
        </div>
        {returnRepo()}
      </div>
    </div>
  )
}


const ScrollDownBtn = () => {
  const portfolio = document.getElementById('portfolio');

  function handleScrollDown() {
    if (!portfolio)
      return;
    portfolio.scrollIntoView({ behavior: 'smooth' });
  }

  return (
    <button id="scroll-down-btn" onClick={handleScrollDown}>
      <Image
        src="/images/chevron-down.svg"
        alt="Chevron down"
        width={24}
        height={24}
      />
    </button>
  )

}


export default function Home() {

  return (
    <div id="app">
      <section id="header">
        <div id="header-content">
          <div id="title">
            <h1>Clément Boillot</h1>
            <p>18 yo French student in computer science</p>
          </div>
          <div id="contact">
            <a href="https://github.com/drawbu" id="github">
              <Image
                src="/images/github.svg"
                alt="GitHub logo"
                width={24}
                height={24}
              />
            </a>
            <a href="https://twitter.com/drawbu" id="twitter">
              <Image
                src="/images/twitter.svg"
                alt="Twitter logo"
                width={24}
                height={24}
              />
            </a>
            <a href="mailto:clement2104.boillot@gmail.com" id="mail">
              <Image
                src="/images/mail.svg"
                alt="Mail icon"
                width={24}
                height={24}
              />
            </a>
          </div>
        </div>
        <div className="separator"></div>
        <ScrollDownBtn />
      </section>

      <section id="portfolio">
        <div className="line">
          <Card
            bg_src="https://cdn.discordapp.com/attachments/837615989830975519/1093939921733042247/random.png"
            name="Cubes fight"
            link="https://cubes.drawbu.dev"
            repo="https://github.com/drawbu/cubes-fight"
          />
          <Card
            bg_src="https://cdn.discordapp.com/attachments/837615989830975519/1093939921733042247/random.png"
            name="ASCII Art generator"
            link="https://ascii.drawbu.dev"
            repo="https://github.com/drawbu/ASCII-Art-generator"
          />
        </div>
        <div className="line">
          <Card
            bg_src="https://cdn.discordapp.com/attachments/837615989830975519/1093939921733042247/random.png"
            name="My dotfiles"
            repo="https://github.com/drawbu/dotfiles"
          />
          <Card
            bg_src="https://cdn.discordapp.com/attachments/837615989830975519/1093939921733042247/random.png"
            name="Nextbus"
            repo="https://github.com/drawbu/nextbus"
          />
        </div>
        <div className="line">
          <Card
            bg_src="https://cdn.discordapp.com/attachments/837615989830975519/1093939921733042247/random.png"
            name="Nuit de l'info 2022"
            link="http://sex-info.ml/"
            repo="https://github.com/Theo-Morin/NDI-2022"
          />
        </div>
      </section>
    </div>
)
}

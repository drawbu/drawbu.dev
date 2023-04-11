import './globals.css'
import React from 'react';

export const metadata = {
  title: 'Clément Boillot',
  description: 'Personal website and portfolio of Clément Boillot, alias drawbu',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}

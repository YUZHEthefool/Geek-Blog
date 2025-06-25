import React from 'react';
import { Link } from 'react-router-dom';
import '../styles/Header.css';

function Header() {
  return (
    <header className="header">
      <div className="header-container">
        <Link to="/" className="logo">
          <span className="logo-bracket">[</span>
          <span className="logo-text">GeekBlog</span>
          <span className="logo-bracket">]</span>
          <span className="cursor">_</span>
        </Link>
        <nav className="nav">
          <Link to="/" className="nav-link">~/posts</Link>
          <Link to="/create" className="nav-link">~/create</Link>
          <a href="https://github.com" className="nav-link" target="_blank" rel="noopener noreferrer">
            ~/github
          </a>
        </nav>
      </div>
    </header>
  );
}

export default Header;
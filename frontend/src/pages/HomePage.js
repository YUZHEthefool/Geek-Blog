import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getPosts } from '../services/api';
import '../styles/HomePage.css';

function HomePage() {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    fetchPosts();
  }, [page]);

  const fetchPosts = async () => {
    try {
      setLoading(true);
      const data = await getPosts(page);
      setPosts(data.posts || []);
      setTotal(data.total || 0);
    } catch (error) {
      console.error('Failed to fetch posts:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="loading-text">Loading<span className="dots">...</span></div>
      </div>
    );
  }

  return (
    <div className="home-page">
      <div className="terminal-header">
        <span className="terminal-prompt">$</span> ls -la /var/log/posts/
      </div>
      
      <div className="posts-list">
        {posts.map((post) => (
          <article key={post.id} className="post-item">
            <Link to={`/post/${post.id}`} className="post-link">
              <h2 className="post-title">
                <span className="hash">#</span> {post.title}
              </h2>
            </Link>
            <div className="post-meta">
              <span className="meta-item">
                <span className="meta-label">author:</span> {post.author}
              </span>
              <span className="meta-separator">|</span>
              <span className="meta-item">
                <span className="meta-label">date:</span> {formatDate(post.created_at)}
              </span>
              <span className="meta-separator">|</span>
              <span className="meta-item">
                <span className="meta-label">views:</span> {post.view_count}
              </span>
            </div>
            {post.tags && post.tags.length > 0 && (
              <div className="post-tags">
                {post.tags.map((tag, index) => (
                  <span key={index} className="tag">#{tag}</span>
                ))}
              </div>
            )}
          </article>
        ))}
      </div>

      {posts.length === 0 && (
        <div className="no-posts">
          <p>No posts found. Start writing!</p>
        </div>
      )}

      <div className="pagination">
        <button 
          onClick={() => setPage(page - 1)} 
          disabled={page === 1}
          className="page-btn"
        >
          &lt; prev
        </button>
        <span className="page-info">Page {page}</span>
        <button 
          onClick={() => setPage(page + 1)} 
          disabled={posts.length < 10}
          className="page-btn"
        >
          next &gt;
        </button>
      </div>
    </div>
  );
}

export default HomePage;
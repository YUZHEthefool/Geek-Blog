import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';
import rehypeHighlight from 'rehype-highlight';
import { getPost, getComments, createComment } from '../services/api';
import CommentSection from '../components/CommentSection';
import '../styles/PostPage.css';
import 'katex/dist/katex.min.css';
import 'highlight.js/styles/github-dark.css';

function PostPage() {
  const { id } = useParams();
  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchPost();
    fetchComments();
  }, [id]);

  const fetchPost = async () => {
    try {
      const data = await getPost(id);
      setPost(data);
    } catch (error) {
      console.error('Failed to fetch post:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchComments = async () => {
    try {
      const data = await getComments(id);
      setComments(data || []);
    } catch (error) {
      console.error('Failed to fetch comments:', error);
    }
  };

  const handleCommentSubmit = async (commentData) => {
    try {
      await createComment({ ...commentData, post_id: id });
      fetchComments();
    } catch (error) {
      console.error('Failed to create comment:', error);
    }
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="loading-text">Loading<span className="dots">...</span></div>
      </div>
    );
  }

  if (!post) {
    return <div className="error">Post not found</div>;
  }

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  return (
    <article className="post-page">
      <header className="post-header">
        <h1 className="post-title">{post.title}</h1>
        <div className="post-meta">
          <span className="meta-item">
            <span className="meta-label">By</span> {post.author}
          </span>
          <span className="meta-separator">•</span>
          <span className="meta-item">{formatDate(post.created_at)}</span>
          <span className="meta-separator">•</span>
          <span className="meta-item">{post.view_count} views</span>
        </div>
        {post.tags && post.tags.length > 0 && (
          <div className="post-tags">
            {post.tags.map((tag, index) => (
              <span key={index} className="tag">#{tag}</span>
            ))}
          </div>
        )}
      </header>

      <div className="post-content">
        <ReactMarkdown
          remarkPlugins={[remarkGfm, remarkMath]}
          rehypePlugins={[rehypeKatex, rehypeHighlight]}
        >
          {post.content}
        </ReactMarkdown>
      </div>

      <CommentSection
        comments={comments}
        onSubmit={handleCommentSubmit}
      />
    </article>
  );
}

export default PostPage;
import React, { useState } from 'react';
import '../styles/CommentSection.css';

function CommentSection({ comments, onSubmit }) {
  const [author, setAuthor] = useState('');
  const [email, setEmail] = useState('');
  const [content, setContent] = useState('');
  const [replyTo, setReplyTo] = useState(null);

  const handleSubmit = (e) => {
    e.preventDefault();
    
    const commentData = {
      author,
      email,
      content,
      parent_id: replyTo
    };
    
    onSubmit(commentData);
    
    // Reset form
    setAuthor('');
    setEmail('');
    setContent('');
    setReplyTo(null);
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const renderComments = (parentId = null, level = 0) => {
    return comments
      .filter(comment => comment.parent_id === parentId)
      .map(comment => (
        <div key={comment.id} className={`comment ${level > 0 ? 'comment-reply' : ''}`}>
          <div className="comment-header">
            <span className="comment-author">{comment.author}</span>
            <span className="comment-date">{formatDate(comment.created_at)}</span>
          </div>
          <div className="comment-content">{comment.content}</div>
          <button 
            className="reply-btn"
            onClick={() => setReplyTo(comment.id)}
          >
            Reply
          </button>
          {renderComments(comment.id, level + 1)}
        </div>
      ));
  };

  return (
    <section className="comment-section">
      <h2>Comments ({comments.length})</h2>
      
      <form onSubmit={handleSubmit} className="comment-form">
        <h3>{replyTo ? 'Reply to comment' : 'Leave a comment'}</h3>
        {replyTo && (
          <button 
            type="button" 
            onClick={() => setReplyTo(null)}
            className="cancel-reply"
          >
            Cancel reply
          </button>
        )}
        
        <div className="form-row">
          <input
            type="text"
            placeholder="Name"
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            required
          />
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        
        <textarea
          placeholder="Write your comment..."
          value={content}
          onChange={(e) => setContent(e.target.value)}
          required
          rows="4"
        />
        
        <button type="submit" className="submit-btn">Post Comment</button>
      </form>
      
      <div className="comments-list">
        {renderComments()}
      </div>
    </section>
  );
}

export default CommentSection;
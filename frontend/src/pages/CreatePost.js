import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';
import rehypeHighlight from 'rehype-highlight';
import { createPost, uploadImage } from '../services/api';
import '../styles/CreatePost.css';

function CreatePost() {
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [author, setAuthor] = useState('');
  const [tags, setTags] = useState('');
  const [preview, setPreview] = useState(false);
  const [uploading, setUploading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    try {
      const postData = {
        title,
        content,
        author,
        tags: tags.split(',').map(tag => tag.trim()).filter(tag => tag)
      };
      
      const response = await createPost(postData);
      navigate(`/post/${response.id}`);
    } catch (error) {
      console.error('Failed to create post:', error);
      alert('Failed to create post');
    }
  };

  const handleImageUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    setUploading(true);
    try {
      const response = await uploadImage(file);
      const imageMarkdown = `![${file.name}](${response.url})`;
      setContent(content + '\n' + imageMarkdown);
    } catch (error) {
      console.error('Failed to upload image:', error);
      alert('Failed to upload image');
    } finally {
      setUploading(false);
    }
  };

  const insertMarkdown = (before, after = '') => {
    const textarea = document.getElementById('content-textarea');
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const selectedText = content.substring(start, end);
    const newText = content.substring(0, start) + before + selectedText + after + content.substring(end);
    setContent(newText);
    
    // Restore cursor position
    setTimeout(() => {
      textarea.focus();
      textarea.setSelectionRange(start + before.length, start + before.length + selectedText.length);
    }, 0);
  };

  return (
    <div className="create-post">
      <div className="editor-header">
        <h1>Create New Post</h1>
        <button 
          type="button" 
          onClick={() => setPreview(!preview)}
          className="preview-toggle"
        >
          {preview ? 'Edit' : 'Preview'}
        </button>
      </div>

      {!preview ? (
        <form onSubmit={handleSubmit} className="post-form">
          <div className="form-group">
            <label htmlFor="title">Title</label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              placeholder="Enter post title..."
            />
          </div>

          <div className="form-group">
            <label htmlFor="author">Author</label>
            <input
              type="text"
              id="author"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              required
              placeholder="Your name..."
            />
          </div>

          <div className="form-group">
            <label htmlFor="tags">Tags (comma separated)</label>
            <input
              type="text"
              id="tags"
              value={tags}
              onChange={(e) => setTags(e.target.value)}
              placeholder="javascript, react, nodejs..."
            />
          </div>

          <div className="form-group">
            <label htmlFor="content">Content (Markdown + LaTeX)</label>
            <div className="editor-toolbar">
              <button type="button" onClick={() => insertMarkdown('**', '**')} title="Bold">B</button>
              <button type="button" onClick={() => insertMarkdown('*', '*')} title="Italic">I</button>
              <button type="button" onClick={() => insertMarkdown('`', '`')} title="Code">{'<>'}</button>
              <button type="button" onClick={() => insertMarkdown('\n```\n', '\n```\n')} title="Code Block">{'{ }'}</button>
              <button type="button" onClick={() => insertMarkdown('$', '$')} title="Inline Math">∑</button>
              <button type="button" onClick={() => insertMarkdown('\n$$\n', '\n$$\n')} title="Block Math">∫</button>
              <button type="button" onClick={() => insertMarkdown('[', '](url)')} title="Link">🔗</button>
              <label className="upload-btn" title="Upload Image">
                📷
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleImageUpload}
                  style={{ display: 'none' }}
                />
              </label>
              {uploading && <span className="uploading">Uploading...</span>}
            </div>
            <textarea
              id="content-textarea"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              required
              rows="20"
              placeholder="Write your post in Markdown format..."
            />
          </div>

          <div className="form-actions">
            <button type="submit" className="submit-btn">Publish Post</button>
            <button type="button" onClick={() => navigate('/')} className="cancel-btn">Cancel</button>
          </div>
        </form>
      ) : (
        <div className="preview-container">
          <article className="post-preview">
            <h1>{title || 'Untitled'}</h1>
            <div className="post-meta">
              <span>By {author || 'Anonymous'}</span>
              <span> • </span>
              <span>{new Date().toLocaleDateString()}</span>
            </div>
            {tags && (
              <div className="post-tags">
                {tags.split(',').map((tag, index) => (
                  <span key={index} className="tag">#{tag.trim()}</span>
                ))}
              </div>
            )}
            <div className="post-content">
              <ReactMarkdown
                remarkPlugins={[remarkGfm, remarkMath]}
                rehypePlugins={[rehypeKatex, rehypeHighlight]}
              >
                {content || '*No content yet*'}
              </ReactMarkdown>
            </div>
          </article>
        </div>
      )}
    </div>
  );
}

export default CreatePost;
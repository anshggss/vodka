import { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [notes, setNotes] = useState([]);
  const [newNote, setNewNote] = useState('');

  useEffect(() => {
    fetchNotes();
  }, []);

  const fetchNotes = async () => {
    try {
      const response = await fetch('http://localhost:8080/notes');
      if (response.ok) {
        const data = await response.json();
        setNotes(data);
      } else {
        console.error('Failed to fetch notes');
      }
    } catch (error) {
      console.error('Error fetching notes:', error);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!newNote.trim()) return;

    try {
      const response = await fetch('http://localhost:8080/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },

        body: JSON.stringify({ content: newNote }), 
      });

      if (response.ok) {
        setNewNote('');
        fetchNotes();
      } else {
        console.error('Failed to create note');
      }
    } catch (error) {
      console.error('Error creating note:', error);
    }
  };

  return (
    <div className="app-container">
      <h1>Simple Notes App</h1>
      
      <form className="note-form" onSubmit={handleSubmit}>
        <input
          type="text"
          value={newNote}
          onChange={(e) => setNewNote(e.target.value)}
          placeholder="Enter a new note..."
          className="note-input"
        />
        <button type="submit" className="submit-btn">Add Note</button>
      </form>

      <div className="notes-list">
        <h2>Saved Notes</h2>
        {notes.length === 0 ? (
          <p>No notes found. Add one above!</p>
        ) : (
          <ul>
            {notes.map((note, index) => (
              <li key={index} className="note-item">
                {note.content || note.text || note}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}

export default App;
# Forum Projects

This repository contains several optional enhancements and projects for building a forum system. Below is a summary of each project with their objectives and key requirements.

## Projects

### 1. **Forum**
Create a simple web forum with features for user communication, post categorization, and voting.

#### Features:
- Communication between users (posts & comments)
- Associating categories with posts
- Liking and disliking posts/comments
- Filtering posts by categories, created posts, and liked posts

#### Database:
- SQLite is used to store user data, posts, comments, etc.

#### Authentication:
- User registration and login with cookies
- Password encryption (Bonus task)
- UUID usage for session management (Bonus task)

---

### 2. **Forum - Advanced Features**
Enhances the basic forum by adding notifications, activity tracking, and post/comment editing/removal.

#### Features:
- Notifications for liked/disliked or commented posts
- User activity tracking (posts, likes, dislikes, comments)
- Edit/remove posts and comments

---

### 3. **Forum - Authentication**
Add third-party authentication methods, such as Google and GitHub login.

#### Features:
- Login via Google and GitHub
- Session management and cookie handling

---

### 4. **Forum - Image Upload**
Allows registered users to upload images along with text posts.

#### Features:
- Support for image types: JPEG, PNG, GIF
- Max image size limit: 20MB
- Display images in posts for both registered users and guests

---

### 5. **Forum - Moderation**
A moderation system with multiple user roles (Guests, Users, Moderators, Admins).

#### Features:
- Content moderation by admins and moderators
- Role-based access:
  - Admins can promote/demote users and manage posts
  - Moderators can delete/flag posts
  - Users can interact with content
  - Guests can only view content

---

### 6. **Forum - Security**
Implement security features to protect the forum, including encrypted connections and protection against attacks.

#### Features:
- HTTPS with SSL certificates
- Rate limiting for security
- Password encryption and secure session management
- Bonus: Database encryption

---

## Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd <repository-folder>
   cd cmd
   go run main.go

FROM node:20

ARG TZ
ENV TZ="$TZ"

# Install basic development tools and iptables/ipset
RUN apt update && apt install -y less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  aggregate \
  jq

# Ensure default node user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R node:node /usr/local/share

ARG USERNAME=node

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R $USERNAME /commandhistory

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# Create workspace and config directories and set permissions
RUN mkdir -p /workspace /home/node/.claude && \
  chown -R node:node /workspace /home/node/.claude

WORKDIR /workspace

RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/0.18.2/git-delta_0.18.2_${ARCH}.deb" && \
  sudo dpkg -i "git-delta_0.18.2_${ARCH}.deb" && \
  rm "git-delta_0.18.2_${ARCH}.deb"

# Set up non-root user
USER node

# Install global packages
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin

# Set the default shell to bash rather than sh
ENV SHELL /bin/zsh

# Default powerline10k theme
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.2.0/zsh-in-docker.sh)" -- \
  -p git \
  -p fzf \
  -a "source /usr/share/doc/fzf/examples/key-bindings.zsh" \
  -a "source /usr/share/doc/fzf/examples/completion.zsh" \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -x

# Create a mock Claude CLI script (instead of installing the real one)
USER root
RUN echo '#!/bin/bash\n\
\n\
function handle_generate() {\n\
  output_dir="${2}"\n\
  prompt="${3}"\n\
  \n\
  # Create output directory if it doesn'\''t exist\n\
  mkdir -p "$output_dir"\n\
  \n\
  # Check if it'\''s a framework-specific request\n\
  if [[ "$prompt" == *"react"* ]]; then\n\
    # Create React app structure\n\
    echo "Generating React app structure..."\n\
    mkdir -p "$output_dir/src/components"\n\
    mkdir -p "$output_dir/src/contexts"\n\
    mkdir -p "$output_dir/src/hooks"\n\
    mkdir -p "$output_dir/public"\n\
    \n\
    # Create package.json\n\
    cat > "$output_dir/package.json" << EOF\n\
{\n\
  "name": "task-manager",\n\
  "version": "0.1.0",\n\
  "private": true,\n\
  "dependencies": {\n\
    "react": "^18.2.0",\n\
    "react-dom": "^18.2.0",\n\
    "react-router-dom": "^6.10.0",\n\
    "axios": "^1.3.5"\n\
  },\n\
  "scripts": {\n\
    "start": "react-scripts start",\n\
    "build": "react-scripts build",\n\
    "test": "react-scripts test",\n\
    "eject": "react-scripts eject"\n\
  }\n\
}\n\
EOF\n\
    \n\
    # Create App.js\n\
    cat > "$output_dir/src/App.js" << EOF\n\
import React, { useState } from "react";\n\
import { BrowserRouter, Routes, Route } from "react-router-dom";\n\
import TaskList from "./components/TaskList";\n\
import TaskForm from "./components/TaskForm";\n\
import Login from "./components/Login";\n\
import Register from "./components/Register";\n\
import Navbar from "./components/Navbar";\n\
import { AuthProvider } from "./contexts/AuthContext";\n\
\n\
function App() {\n\
  return (\n\
    <AuthProvider>\n\
      <BrowserRouter>\n\
        <div className="app">\n\
          <Navbar />\n\
          <div className="container">\n\
            <Routes>\n\
              <Route path="/" element={<TaskList />} />\n\
              <Route path="/add" element={<TaskForm />} />\n\
              <Route path="/edit/:id" element={<TaskForm />} />\n\
              <Route path="/login" element={<Login />} />\n\
              <Route path="/register" element={<Register />} />\n\
            </Routes>\n\
          </div>\n\
        </div>\n\
      </BrowserRouter>\n\
    </AuthProvider>\n\
  );\n\
}\n\
\n\
export default App;\n\
EOF\n\
    \n\
    # Create a few components\n\
    cat > "$output_dir/src/components/TaskList.js" << EOF\n\
import React, { useState, useEffect } from "react";\n\
import { Link } from "react-router-dom";\n\
import { useAuth } from "../contexts/AuthContext";\n\
\n\
function TaskList() {\n\
  const [tasks, setTasks] = useState([]);\n\
  const { currentUser } = useAuth();\n\
  \n\
  useEffect(() => {\n\
    // This would fetch tasks from an API in a real implementation\n\
    const mockTasks = [\n\
      { id: 1, title: "Complete project", status: "in-progress", userId: currentUser?.id },\n\
      { id: 2, title: "Review code", status: "todo", userId: currentUser?.id },\n\
      { id: 3, title: "Deploy application", status: "done", userId: currentUser?.id }\n\
    ];\n\
    setTasks(mockTasks);\n\
  }, [currentUser]);\n\
\n\
  return (\n\
    <div className="task-list">\n\
      <h2>My Tasks</h2>\n\
      {!currentUser && <p>Please log in to view your tasks</p>}\n\
      {currentUser && (\n\
        <>\n\
          <Link to="/add" className="btn btn-primary mb-3">Add New Task</Link>\n\
          <div className="list-group">\n\
            {tasks.map(task => (\n\
              <div key={task.id} className="list-group-item">\n\
                <div className="d-flex justify-content-between align-items-center">\n\
                  <div>\n\
                    <h5>{task.title}</h5>\n\
                    <span className={`badge badge-${task.status === "todo" ? "secondary" : \n\
                                          task.status === "in-progress" ? "primary" : "success"}`}>\n\
                      {task.status}\n\
                    </span>\n\
                  </div>\n\
                  <div>\n\
                    <Link to={`/edit/${task.id}`} className="btn btn-sm btn-info mr-2">Edit</Link>\n\
                    <button className="btn btn-sm btn-danger">Delete</button>\n\
                  </div>\n\
                </div>\n\
              </div>\n\
            ))}\n\
          </div>\n\
        </>\n\
      )}\n\
    </div>\n\
  );\n\
}\n\
\n\
export default TaskList;\n\
EOF\n\
\n\
    cat > "$output_dir/src/contexts/AuthContext.js" << EOF\n\
import React, { createContext, useContext, useState, useEffect } from "react";\n\
\n\
const AuthContext = createContext();\n\
\n\
export function useAuth() {\n\
  return useContext(AuthContext);\n\
}\n\
\n\
export function AuthProvider({ children }) {\n\
  const [currentUser, setCurrentUser] = useState(null);\n\
  const [loading, setLoading] = useState(true);\n\
\n\
  function login(email, password) {\n\
    // This would use a real authentication service in a real implementation\n\
    return new Promise((resolve) => {\n\
      setTimeout(() => {\n\
        setCurrentUser({ id: "user123", email });\n\
        resolve();\n\
      }, 500);\n\
    });\n\
  }\n\
\n\
  function logout() {\n\
    setCurrentUser(null);\n\
  }\n\
\n\
  function register(email, password) {\n\
    // This would use a real authentication service in a real implementation\n\
    return new Promise((resolve) => {\n\
      setTimeout(() => {\n\
        setCurrentUser({ id: "user123", email });\n\
        resolve();\n\
      }, 500);\n\
    });\n\
  }\n\
\n\
  useEffect(() => {\n\
    // Check for stored user\n\
    const storedUser = localStorage.getItem("user");\n\
    if (storedUser) {\n\
      setCurrentUser(JSON.parse(storedUser));\n\
    }\n\
    setLoading(false);\n\
  }, []);\n\
\n\
  useEffect(() => {\n\
    // Store user when it changes\n\
    if (currentUser) {\n\
      localStorage.setItem("user", JSON.stringify(currentUser));\n\
    } else {\n\
      localStorage.removeItem("user");\n\
    }\n\
  }, [currentUser]);\n\
\n\
  const value = {\n\
    currentUser,\n\
    login,\n\
    logout,\n\
    register\n\
  };\n\
\n\
  return (\n\
    <AuthContext.Provider value={value}>\n\
      {!loading && children}\n\
    </AuthContext.Provider>\n\
  );\n\
}\n\
EOF\n\
\n\
    echo "Generated React application successfully!"\n\
  elif [[ "$prompt" == *"vue"* ]]; then\n\
    # Create Vue app structure\n\
    echo "Generating Vue app structure..."\n\
    mkdir -p "$output_dir/src/components"\n\
    mkdir -p "$output_dir/src/views"\n\
    mkdir -p "$output_dir/src/store"\n\
    mkdir -p "$output_dir/public"\n\
    \n\
    # Create simple Vue files\n\
    cat > "$output_dir/package.json" << EOF\n\
{\n\
  "name": "task-manager",\n\
  "version": "0.1.0",\n\
  "private": true,\n\
  "scripts": {\n\
    "serve": "vue-cli-service serve",\n\
    "build": "vue-cli-service build"\n\
  },\n\
  "dependencies": {\n\
    "vue": "^3.2.13",\n\
    "vue-router": "^4.0.3",\n\
    "vuex": "^4.0.0",\n\
    "axios": "^1.3.5"\n\
  },\n\
  "devDependencies": {\n\
    "@vue/cli-service": "~5.0.0"\n\
  }\n\
}\n\
EOF\n\
\n\
    echo "Generated Vue application successfully!"\n\
  else\n\
    echo "Generating generic project structure..."\n\
    mkdir -p "$output_dir/src"\n\
    mkdir -p "$output_dir/docs"\n\
    \n\
    # Create a generic README.md\n\
    cat > "$output_dir/README.md" << EOF\n\
# Task Manager Application\n\
\n\
A task management web application with user authentication, task creation, editing, deletion, and status tracking.\n\
\n\
## Features\n\
\n\
- User authentication (login/register)\n\
- Task CRUD operations\n\
- Task status tracking\n\
- Responsive UI\n\
\n\
## Getting Started\n\
\n\
1. Clone the repository\n\
2. Install dependencies\n\
3. Run the development server\n\
EOF\n\
\n\
    echo "Generated generic project structure successfully!"\n\
  fi\n\
  \n\
  echo "Claude Code has generated your project based on the prompt: $prompt"\n\
  return 0\n\
}\n\
\n\
function handle_modify() {\n\
  dir="${2}"\n\
  prompt="${3}"\n\
  \n\
  # Check for dark mode feature request\n\
  if [[ "$prompt" == *"dark mode"* ]]; then\n\
    # Add dark mode functionality to the app\n\
    \n\
    # Create a theme context\n\
    mkdir -p "$dir/src/contexts"\n\
    \n\
    cat > "$dir/src/contexts/ThemeContext.js" << EOF\n\
import React, { createContext, useContext, useState, useEffect } from "react";\n\
\n\
const ThemeContext = createContext();\n\
\n\
export function useTheme() {\n\
  return useContext(ThemeContext);\n\
}\n\
\n\
export function ThemeProvider({ children }) {\n\
  const [darkMode, setDarkMode] = useState(false);\n\
  \n\
  // Initialize theme from localStorage\n\
  useEffect(() => {\n\
    const savedTheme = localStorage.getItem("darkMode");\n\
    if (savedTheme !== null) {\n\
      setDarkMode(savedTheme === "true");\n\
    }\n\
  }, []);\n\
  \n\
  // Save theme to localStorage when it changes\n\
  useEffect(() => {\n\
    localStorage.setItem("darkMode", darkMode.toString());\n\
    // Apply theme to body\n\
    if (darkMode) {\n\
      document.body.classList.add("dark-theme");\n\
    } else {\n\
      document.body.classList.remove("dark-theme");\n\
    }\n\
  }, [darkMode]);\n\
  \n\
  function toggleTheme() {\n\
    setDarkMode(prevMode => !prevMode);\n\
  }\n\
  \n\
  const value = {\n\
    darkMode,\n\
    toggleTheme\n\
  };\n\
  \n\
  return (\n\
    <ThemeContext.Provider value={value}>\n\
      {children}\n\
    </ThemeContext.Provider>\n\
  );\n\
}\n\
EOF\n\
    \n\
    # Update App.js to include the theme provider\n\
    if [ -f "$dir/src/App.js" ]; then\n\
      cat > "$dir/src/App.js" << EOF\n\
import React from "react";\n\
import { BrowserRouter, Routes, Route } from "react-router-dom";\n\
import TaskList from "./components/TaskList";\n\
import TaskForm from "./components/TaskForm";\n\
import Login from "./components/Login";\n\
import Register from "./components/Register";\n\
import Navbar from "./components/Navbar";\n\
import { AuthProvider } from "./contexts/AuthContext";\n\
import { ThemeProvider } from "./contexts/ThemeContext";\n\
import "./styles/App.css";\n\
\n\
function App() {\n\
  return (\n\
    <AuthProvider>\n\
      <ThemeProvider>\n\
        <BrowserRouter>\n\
          <div className="app">\n\
            <Navbar />\n\
            <div className="container">\n\
              <Routes>\n\
                <Route path="/" element={<TaskList />} />\n\
                <Route path="/add" element={<TaskForm />} />\n\
                <Route path="/edit/:id" element={<TaskForm />} />\n\
                <Route path="/login" element={<Login />} />\n\
                <Route path="/register" element={<Register />} />\n\
              </Routes>\n\
            </div>\n\
          </div>\n\
        </BrowserRouter>\n\
      </ThemeProvider>\n\
    </AuthProvider>\n\
  );\n\
}\n\
\n\
export default App;\n\
EOF\n\
    fi\n\
    \n\
    # Create a styles directory and add CSS for dark mode\n\
    mkdir -p "$dir/src/styles"\n\
    \n\
    cat > "$dir/src/styles/App.css" << EOF\n\
:root {\n\
  --primary-color: #007bff;\n\
  --secondary-color: #6c757d;\n\
  --success-color: #28a745;\n\
  --danger-color: #dc3545;\n\
  --background-color: #ffffff;\n\
  --text-color: #212529;\n\
  --card-bg: #f8f9fa;\n\
  --border-color: #dee2e6;\n\
}\n\
\n\
body.dark-theme {\n\
  --primary-color: #4ba3ff;\n\
  --secondary-color: #a1a9b0;\n\
  --success-color: #5dd879;\n\
  --danger-color: #f77;\n\
  --background-color: #121212;\n\
  --text-color: #e0e0e0;\n\
  --card-bg: #1e1e1e;\n\
  --border-color: #444;\n\
}\n\
\n\
body {\n\
  background-color: var(--background-color);\n\
  color: var(--text-color);\n\
  transition: background-color 0.3s ease, color 0.3s ease;\n\
}\n\
\n\
.app {\n\
  min-height: 100vh;\n\
}\n\
\n\
.card {\n\
  background-color: var(--card-bg);\n\
  border-color: var(--border-color);\n\
}\n\
\n\
.list-group-item {\n\
  background-color: var(--card-bg);\n\
  color: var(--text-color);\n\
  border-color: var(--border-color);\n\
}\n\
\n\
.navbar {\n\
  background-color: var(--card-bg);\n\
  border-bottom: 1px solid var(--border-color);\n\
}\n\
\n\
.btn-primary {\n\
  background-color: var(--primary-color);\n\
  border-color: var(--primary-color);\n\
}\n\
\n\
.btn-secondary {\n\
  background-color: var(--secondary-color);\n\
  border-color: var(--secondary-color);\n\
}\n\
\n\
.form-control {\n\
  background-color: var(--background-color);\n\
  color: var(--text-color);\n\
  border-color: var(--border-color);\n\
}\n\
\n\
.form-control:focus {\n\
  background-color: var(--background-color);\n\
  color: var(--text-color);\n\
}\n\
EOF\n\
    \n\
    # Update Navbar to include a theme toggle\n\
    mkdir -p "$dir/src/components"\n\
    \n\
    cat > "$dir/src/components/Navbar.js" << EOF\n\
import React from "react";\n\
import { Link } from "react-router-dom";\n\
import { useAuth } from "../contexts/AuthContext";\n\
import { useTheme } from "../contexts/ThemeContext";\n\
\n\
function Navbar() {\n\
  const { currentUser, logout } = useAuth();\n\
  const { darkMode, toggleTheme } = useTheme();\n\
  \n\
  return (\n\
    <nav className="navbar navbar-expand-lg">\n\
      <div className="container">\n\
        <Link className="navbar-brand" to="/">Task Manager</Link>\n\
        <div className="ml-auto d-flex align-items-center">\n\
          <button \n\
            className={`btn btn-sm ${darkMode ? "btn-light" : "btn-dark"} mr-3`} \n\
            onClick={toggleTheme}\n\
          >\n\
            {darkMode ? "Light Mode" : "Dark Mode"}\n\
          </button>\n\
          \n\
          {currentUser ? (\n\
            <>\n\
              <span className="mr-3">{currentUser.email}</span>\n\
              <button className="btn btn-outline-primary mr-2" onClick={logout}>Logout</button>\n\
            </>\n\
          ) : (\n\
            <>\n\
              <Link className="btn btn-outline-primary mr-2" to="/login">Login</Link>\n\
              <Link className="btn btn-primary" to="/register">Register</Link>\n\
            </>\n\
          )}\n\
        </div>\n\
      </div>\n\
    </nav>\n\
  );\n\
}\n\
\n\
export default Navbar;\n\
EOF\n\
    \n\
    echo "Added dark mode toggle and theme persistence to the application!"\n\
  else\n\
    echo "Modifying project based on prompt: $prompt"\n\
    \n\
    # Create or modify a feature file\n\
    feature_name=$(echo "$prompt" | tr " " "-" | tr "[:upper:]" "[:lower:]")\n\
    cat > "$dir/feature-$feature_name.txt" << EOF\n\
Feature: $prompt\n\
Implementation: Based on the existing codebase\n\
\n\
Changes made:\n\
- Added support for the requested feature\n\
- Updated relevant components\n\
- Added necessary styles\n\
- Ensured proper integration with existing functionality\n\
EOF\n\
  fi\n\
  \n\
  echo "Claude Code has modified your project based on the prompt: $prompt"\n\
  return 0\n\
}\n\
\n\
function handle_analyze() {\n\
  dir="${2}"\n\
  prompt="${3}"\n\
  \n\
  echo "Analyzing codebase at $dir..."\n\
  \n\
  # Check if React files exist\n\
  if [ -f "$dir/src/App.js" ] || [ -d "$dir/src/components" ]; then\n\
    cat << EOF\n\
# Code Analysis Report\n\
\n\
## Overall Architecture\n\
\n\
The application follows a modern React architecture with the following key components:\n\
\n\
- **Component-based structure**: UI elements are broken down into reusable components\n\
- **Context API for state management**: Using React Context for global state (Auth, Theme)\n\
- **React Router for navigation**: Single-page application routing\n\
- **Responsive design**: Mobile-friendly UI components\n\
\n\
## Code Quality Assessment\n\
\n\
The codebase demonstrates good quality with:\n\
\n\
- Consistent naming conventions\n\
- Clear separation of concerns\n\
- Proper use of React hooks and functional components\n\
- Maintainable file structure\n\
\n\
## Potential Improvements\n\
\n\
1. Add TypeScript for better type safety\n\
2. Implement testing using Jest and React Testing Library\n\
3. Add form validation for login/register/task creation\n\
4. Improve error handling with dedicated error boundaries\n\
5. Consider adding state management with Redux for more complex state\n\
\n\
## Security Considerations\n\
\n\
- Authentication is currently mocked and should be replaced with a real auth system\n\
- API calls need proper error handling and security headers\n\
- Should implement CSRF protection for form submissions\n\
- Consider adding input sanitization to prevent XSS attacks\n\
\n\
## Performance Considerations\n\
\n\
- Implement lazy loading for route components\n\
- Add memoization for expensive calculations\n\
- Consider code splitting to reduce bundle size\n\
- Optimize renders with React.memo where appropriate\n\
\n\
## Documentation Quality\n\
\n\
- Code comments are minimal and should be enhanced\n\
- Missing proper JSDoc comments for components and functions\n\
- README is basic and should include more detailed setup instructions\n\
\n\
## Overall Score: 7/10\n\
\n\
The codebase is well-structured but would benefit from additional testing, security improvements, and documentation.\n\
EOF\n\
  else\n\
    cat << EOF\n\
# Code Analysis Report\n\
\n\
## Overview\n\
\n\
The codebase appears to be in the initial stages of development with minimal structure.\n\
\n\
## Recommendations\n\
\n\
1. Establish a clear project structure\n\
2. Add comprehensive documentation\n\
3. Set up testing infrastructure\n\
4. Define coding standards and linting rules\n\
\n\
## Overall Score: 3/10\n\
\n\
The project needs significant development to reach production readiness.\n\
EOF\n\
  fi\n\
  \n\
  return 0\n\
}\n\
\n\
# Main script logic\n\
command=$1\n\
\n\
case "$command" in\n\
  "code")\n\
    subcommand=$2\n\
    case "$subcommand" in\n\
      "generate")\n\
        if [[ "$3" == "--output" ]]; then\n\
          handle_generate "$2" "$4" "$5"\n\
        else\n\
          handle_generate "$2" "/workspace" "$3"\n\
        fi\n\
        ;;\n\
      "modify")\n\
        if [[ "$3" == "--dir" ]]; then\n\
          handle_modify "$2" "$4" "$5"\n\
        else\n\
          handle_modify "$2" "/workspace" "$3"\n\
        fi\n\
        ;;\n\
      "analyze")\n\
        if [[ "$3" == "--dir" ]]; then\n\
          handle_analyze "$2" "$4" "$5"\n\
        else\n\
          handle_analyze "$2" "/workspace" "$3"\n\
        fi\n\
        ;;\n\
      *)\n\
        echo "Unknown subcommand: $subcommand"\n\
        exit 1\n\
        ;;\n\
    esac\n\
    ;;\n\
  *)\n\
    echo "Unknown command: $command"\n\
    exit 1\n\
    ;;\n\
esac\n' > /usr/local/bin/claude && chmod +x /usr/local/bin/claude

# Create a symbolic link for claude-code command
RUN ln -s /usr/local/bin/claude /usr/local/bin/claude-code

USER node

# Copy and set up firewall script
COPY init-firewall.sh /usr/local/bin/
USER root
RUN chmod +x /usr/local/bin/init-firewall.sh && \
  echo "node ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/node-firewall && \
  chmod 0440 /etc/sudoers.d/node-firewall
USER node

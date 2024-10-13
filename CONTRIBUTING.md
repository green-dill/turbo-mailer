# Contributing Guidelines

Thank you for considering contributing to our project! We welcome all forms of contributions, whether it's new features, bug fixes, or documentation improvements. Please follow these guidelines to ensure your contributions can be smoothly accepted.

## How to Contribute

1. Clone the repository
2. Create your branch using the appropriate naming convention:
   - For bug fixes: `git checkout -b bugfix/DescriptiveBugfixName`
   - For new features: `git checkout -b feature/DescriptiveFeatureName`
   - For hotfixes: `git checkout -b hotfix/DescriptiveHotfixName`
3. Make your changes in the new branch
4. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
5. Push to the branch (`git push origin feature/AmazingFeature`)
6. Open a Pull Request

Note: Do not commit directly to `master` or `develop` branches. Always create a new branch for your work.

## Code Style

- Ensure your code adheres to the existing style of the project.
- Use meaningful variable and function names.
- Add necessary comments, but avoid over-commenting.
- Follow the language-specific style guide if one exists for the project.

## Commit Message Guidelines

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 50 characters or less
- Reference issues and pull requests liberally in the description
- Start the commit message with an applicable type:
    * `chore:` for routine tasks or maintenance
    * `docs:` for documentation updates
    * `style:` for code style changes (formatting, missing semi-colons, etc)
    * `refactor:` for code refactoring
    * `perf:` for performance improvements
    * `test:` for adding or modifying tests
    * `bugfix:` or `fix:` for bug fixes
    * `improvement:` for general improvements

Example: `docs: Update README with new installation instructions`

## Reporting Bugs

When reporting bugs, please include:

- Your operating system name and version
- Any relevant details about your browser if applicable
- Detailed steps to reproduce the issue
- What you expected to happen
- What actually happened

## Suggesting Enhancements

If you have ideas for improving the project, we'd love to hear them. Please open an issue and provide the following information:

- A clear and concise description of your idea
- Explain why this feature would be useful to most users
- Suggest some ways to implement it (if possible)
- Describe alternatives you've considered

## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a build.
2. Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations, and container parameters.
3. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent.
4. Your Pull Request will be reviewed by at least two other developers. They may request changes or ask for clarifications.
5. Once approved, your Pull Request will be merged by a project maintainer.

Thank you again for your contributions!

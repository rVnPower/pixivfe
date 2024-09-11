---
hide:
  - toc
---

# Developer Documentation

This section contains documentation related to the development of PixivFE. While primarily intended for developers, instance administrators may find this information useful for gaining a deeper understanding of the application and its future roadmap.

## Overview

The developer documentation is comprised of the following files:

1. [Coding Tips](coding-tips.md): Useful tricks and best practices for PixivFE development, including handling of cookies and Jet template specifics.

2. [Design Flaws](design-flaws.md): Documentation of current design issues in PixivFE, both frontend and backend, with potential solutions.

3. [Feature Ideas](feature-ideas.md): Proposals for potential features or redesigns that can be implemented into PixivFE, including Pixivision integration, Sketch support, and UI improvements.

4. [Framework Migration](framework-migration.md): Documentation of the migration from gofiber to net/http. This migration was completed in [release v2.8](https://codeberg.org/VnPower/PixivFE/src/tag/v2.8) (commit [1ac6f40b57](https://codeberg.org/VnPower/PixivFE/src/commit/1ac6f40b57608b1576d9e812698e95958c91c626)).

5. [Helpful Resources](helpful-resources.md): External links to materials and resources that can aid in PixivFE development, including other Pixiv-related projects and tools.

6. [Roadmap](roadmap.md): Planned features and improvements for PixivFE, categorized by implementation status and priority.

7. [Testing](testing.md): Information about the current state of testing in PixivFE and considerations for future testing strategies.

## Features in Development

The Features section contains detailed information about specific features that are currently in development or planned for future implementation:

1. [Caching](features/caching.md): Outlines caching strategies for various types of content, including images, JSON requests, and rendered pages. Also discusses the potential for predictive caching to improve browsing experience.

2. [Novels](features/novels.md): Lists planned improvements for the novel reading experience, including UI enhancements, support for novel series, and additional features like furigana support and vertical text display.

3. [Tracing and Flamegraphs](features/tracing-flamegraph.md): Explains the tracing implemented in PixivFE for both Pixiv website requests and server requests, and how to view flamegraphs for performance analysis.

4. [User Customization](features/user-customization.md): Details potential per-user customization options for various aspects of PixivFE, including site-wide settings, novel reading preferences, artwork filtering, and search options.

## Contributing

Developers interested in contributing to PixivFE are encouraged to review the documentation in this section, particularly [Roadmap](roadmap.md), [Feature Ideas](feature-ideas.md), and the Features section. These resources will help identify areas where meaningful contributions can be made.

For any questions or discussions related to PixivFE development, please refer to the [project's issue tracker on Codeberg](https://codeberg.org/VnPower/PixivFE/issues).

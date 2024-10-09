vcl 4.1;

backend pixivfe {
  .host = "pixivfe";
  .port = "8282";
}

sub vcl_recv {
  if (req.http.host == "caddy.localhost") {
    set req.backend_hint = pixivfe;

    # Step 1: Check for login status
    if (req.http.cookie ~ "pixivfe-Token") {
      set req.http.X-User-Logged-In = "true";
    } else {
      set req.http.X-User-Logged-In = "false";
    }

    # For logged-in users, include the token in the cache key
    if (req.http.X-User-Logged-In == "true") {
      set req.http.X-User-Token = regsub(req.http.cookie, ".*pixivfe-Token=([^;]+).*", "\1");
    }
  }
}

sub vcl_backend_response {
  if (bereq.backend == pixivfe) {
    # Step 2: Implement different caching strategies

    # set grace to 1 minute globally
    set beresp.grace = 1m;

    ## landing page; should be cached conservatively when it contains personalised content for logged in users, but be more aggressive for the shared page shown to anonymous visitors
    if (bereq.url == "/") {
      if (beresp.status == 200) {
        if (bereq.http.X-User-Logged-In == "false") {
          # Non-logged in users: cache for 8h
          set beresp.ttl = 8h;
          set beresp.http.Cache-Control = "public, max-age=28800";
        } else {
          # Logged in users: cache for 30 seconds, vary by user token
          set beresp.ttl = 30s;
          set beresp.uncacheable = false;
          set beresp.http.Cache-Control = "private, max-age=30";
          set beresp.http.Vary = "X-User-Token";
        }
      }
      return (deliver);
    }

    ## pages with user content that changes dynamically every request as per the Pixiv upstream; use a moderate caching strategy, we don't need to fetch new content *every* time
    if (bereq.url ~ "^/discovery" || bereq.url ~ "^/discovery/novel" || bereq.url ~ "^/newest") {
      if (beresp.status == 200) {
        set beresp.ttl = 30s;
        set beresp.uncacheable = false;
        set beresp.http.Cache-Control = "public, max-age=30";
      }
      return (deliver);
    }

    ## pages with user content that changes dynamically, but either over longer timeframes (ranking pages) or not appreciably (loading in under 100ms is probably more important than knowing the *exact* number of bookmarks a work has); we can be more aggressive
    if (bereq.url ~ "^/ranking" || bereq.url ~ "^/rankingCalendar" || bereq.url ~ "^/artworks/" || bereq.url ~ "^/novel/" || bereq.url ~ "^/users/") {
      if (beresp.status == 200) {
        set beresp.ttl = 8h;
        set beresp.http.Cache-Control = "public, max-age=28800";
      }
      return (deliver);
    }

    # proxied static media (images, short video, etc); a aggressive caching strategy is safe
    if (bereq.url ~ "^/proxy/") {
      if (beresp.status == 200) {
        set beresp.ttl = 24h;
        set beresp.http.Cache-Control = "public, max-age=31536000"; # 1 year
      }
      return (deliver);
    }

    # if nothing matches, be safe and dont cache anything
    set beresp.http.Cache-Control = "private, no-store";
  }
}

sub vcl_hash {
  # Step 3: Include user token in hash for logged-in users
  if (req.http.X-User-Logged-In == "true") {
    hash_data(req.http.X-User-Token);
  }
}

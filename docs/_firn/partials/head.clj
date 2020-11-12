(defn head
  [{:keys [build-url site-title site-author site-desc]}]
  [:head
   [:meta {:charset "utf-8"}]
   [:meta {:name "viewport" :content "width=device-width, initial-scale=1.0"}]
   [:meta {:name "author" :content site-author}]
   [:meta {:name "description" :content site-desc}]
   [:title site-title]
   [:script {:async :defer :data-domain "flattrack.gitlab.io" :src "https://plausible.io/js/plausible.js"}]
   [:link {:rel "stylesheet" :href (build-url "/static/css/firn_base.css")}]])
;; <script async defer data-domain="flattrack.gitlab.io" src="https://plausible.io/js/plausible.js"></script>

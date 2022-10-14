package webapi

type PublishedFileDetailsResponse struct {
	Response struct {
		PublishedFileDetails []struct {
			Result                   int    `json:"result"`
			PublishedFileID          string `json:"publishedfileid"`
			Creator                  string `json:"creator"`
			CreatorAppid             int    `json:"creator_appid"`
			ConsumerAppid            int    `json:"consumer_appid"`
			ConsumerShortcutID       int    `json:"consumer_shortcutid"`
			Filename                 string `json:"filename"`
			FileSize                 string `json:"file_size"`
			PreviewFileSize          string `json:"preview_file_size"`
			FileURL                  string `json:"file_url"`
			PreviewURL               string `json:"preview_url"`
			URL                      string `json:"url"`
			HContentFile             string `json:"hcontent_file"`
			HContentPreview          string `json:"hcontent_preview"`
			Title                    string `json:"title"`
			FileDescription          string `json:"file_description"`
			TimeCreated              int    `json:"time_created"`
			TimeUpdated              int    `json:"time_updated"`
			Visibility               int    `json:"visibility"`
			Flags                    int    `json:"flags"`
			WorkshopFile             bool   `json:"workshop_file"`
			WorkshopAccepted         bool   `json:"workshop_accepted"`
			ShowSubscribeAll         bool   `json:"show_subscribe_all"`
			NumCommentsPublic        int    `json:"num_comments_public"`
			Banned                   bool   `json:"banned"`
			BanReason                string `json:"ban_reason"`
			Banner                   string `json:"banner"`
			CanBeDeleted             bool   `json:"can_be_deleted"`
			AppName                  string `json:"app_name"`
			FileType                 int    `json:"file_type"`
			CanSubscribe             bool   `json:"can_subscribe"`
			Subscriptions            int    `json:"subscriptions"`
			Favorited                int    `json:"favorited"`
			Followers                int    `json:"followers"`
			LifetimeSubscriptions    int    `json:"lifetime_subscriptions"`
			LifetimeFavorited        int    `json:"lifetime_favorited"`
			LifetimeFollowers        int    `json:"lifetime_followers"`
			LifetimePlaytime         string `json:"lifetime_playtime"`
			LifetimePlaytimeSessions string `json:"lifetime_playtime_sessions"`
			Views                    int    `json:"views"`
			ImageWidth               int    `json:"image_width,omitempty"`
			ImageHeight              int    `json:"image_height,omitempty"`
			NumChildren              int    `json:"num_children"`
			NumReports               int    `json:"num_reports"`
			Tags                     []struct {
				Tag         string `json:"tag"`
				DisplayName string `json:"display_name"`
			} `json:"tags"`
			Children []struct {
				PublishedFileID string `json:"publishedfileid"`
				SortOrder       int    `json:"sortorder"`
				FileType        int    `json:"file_type"`
			} `json:"children,omitempty"`
			Language                   int    `json:"language"`
			MaybeInappropriateSex      bool   `json:"maybe_inappropriate_sex"`
			MaybeInappropriateViolence bool   `json:"maybe_inappropriate_violence"`
			RevisionChangeNumber       string `json:"revision_change_number"`
			Revision                   int    `json:"revision"`
			BanTextCheckResult         int    `json:"ban_text_check_result"`
		} `json:"publishedfiledetails"`
	} `json:"response"`
}

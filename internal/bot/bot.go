package bot

import (
	"fmt"
	"log/slog"
	"strings"

	discord "github.com/sikigasa/discord-api-go"
	"github.com/sikigasa/todoapp-discordbot/internal/api"
)

// Bot はDiscord Botを表す
type Bot struct {
	client        *discord.Client
	apiClient     *api.Client
	defaultUserID string
	logger        *slog.Logger
}

// New は新しいBotインスタンスを作成する
func New(token string, apiClient *api.Client, defaultUserID string, logger *slog.Logger) (*Bot, error) {
	client := discord.NewClient(token, logger)

	bot := &Bot{
		client:        client,
		apiClient:     apiClient,
		defaultUserID: defaultUserID,
		logger:        logger,
	}

	client.OnInteraction(bot.handleInteraction)

	return bot, nil
}

// Start はBotを起動し、スラッシュコマンドを登録する
func (b *Bot) Start() error {
	if err := b.client.Open(); err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}

	if err := b.registerCommands(); err != nil {
		return fmt.Errorf("failed to register commands: %w", err)
	}

	b.logger.Info("Bot started successfully")
	return nil
}

// Stop はBotを停止する
func (b *Bot) Stop() error {
	b.logger.Info("Shutting down bot...")
	return b.client.Close()
}

// registerCommands はスラッシュコマンドを登録する
func (b *Bot) registerCommands() error {
	commands := []*discord.ApplicationCommand{
		{
			Name:        "todo",
			Description: "TODOを管理する",
			Options: []*discord.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "TODOを作成する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "title",
							Description: "TODOのタイトル",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
						{
							Name:        "description",
							Description: "TODOの説明",
							Type:        discord.OptionTypeString,
							Required:    false,
						},
					},
				},
				{
					Name:        "list",
					Description: "全TODOを一覧表示する",
					Type:        discord.OptionTypeSubCommand,
				},
				{
					Name:        "get",
					Description: "TODOの詳細を取得する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "TODOのID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
					},
				},
				{
					Name:        "complete",
					Description: "TODOを完了にする",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "TODOのID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
					},
				},
				{
					Name:        "update",
					Description: "TODOを更新する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "TODOのID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
						{
							Name:        "title",
							Description: "新しいタイトル",
							Type:        discord.OptionTypeString,
							Required:    false,
						},
						{
							Name:        "description",
							Description: "新しい説明",
							Type:        discord.OptionTypeString,
							Required:    false,
						},
					},
				},
				{
					Name:        "delete",
					Description: "TODOを削除する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "TODOのID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "task",
			Description: "タスクを管理する",
			Options: []*discord.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "タスクを作成する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "project_id",
							Description: "プロジェクトID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
						{
							Name:        "title",
							Description: "タスクのタイトル",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
						{
							Name:        "description",
							Description: "タスクの説明",
							Type:        discord.OptionTypeString,
							Required:    false,
						},
						{
							Name:        "status",
							Description: "ステータス (0: To Do, 1: In Progress, 2: Done)",
							Type:        discord.OptionTypeInteger,
							Required:    false,
							Choices: []*discord.ApplicationCommandOptionChoice{
								{Name: "To Do", Value: 0},
								{Name: "In Progress", Value: 1},
								{Name: "Done", Value: 2},
							},
						},
						{
							Name:        "priority",
							Description: "優先度 (0: Low, 1: Medium, 2: High)",
							Type:        discord.OptionTypeInteger,
							Required:    false,
							Choices: []*discord.ApplicationCommandOptionChoice{
								{Name: "Low", Value: 0},
								{Name: "Medium", Value: 1},
								{Name: "High", Value: 2},
							},
						},
					},
				},
				{
					Name:        "list",
					Description: "プロジェクトのタスク一覧を表示する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "project_id",
							Description: "プロジェクトID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
					},
				},
				{
					Name:        "get",
					Description: "タスクの詳細を取得する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "タスクのID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
					},
				},
				{
					Name:        "update",
					Description: "タスクを更新する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "タスクのID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
						{
							Name:        "title",
							Description: "新しいタイトル",
							Type:        discord.OptionTypeString,
							Required:    false,
						},
						{
							Name:        "description",
							Description: "新しい説明",
							Type:        discord.OptionTypeString,
							Required:    false,
						},
						{
							Name:        "status",
							Description: "ステータス (0: To Do, 1: In Progress, 2: Done)",
							Type:        discord.OptionTypeInteger,
							Required:    false,
							Choices: []*discord.ApplicationCommandOptionChoice{
								{Name: "To Do", Value: 0},
								{Name: "In Progress", Value: 1},
								{Name: "Done", Value: 2},
							},
						},
						{
							Name:        "priority",
							Description: "優先度 (0: Low, 1: Medium, 2: High)",
							Type:        discord.OptionTypeInteger,
							Required:    false,
							Choices: []*discord.ApplicationCommandOptionChoice{
								{Name: "Low", Value: 0},
								{Name: "Medium", Value: 1},
								{Name: "High", Value: 2},
							},
						},
					},
				},
				{
					Name:        "delete",
					Description: "タスクを削除する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "タスクのID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "project",
			Description: "プロジェクトを管理する",
			Options: []*discord.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "プロジェクトを作成する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "title",
							Description: "プロジェクトのタイトル",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
						{
							Name:        "description",
							Description: "プロジェクトの説明",
							Type:        discord.OptionTypeString,
							Required:    false,
						},
					},
				},
				{
					Name:        "list",
					Description: "プロジェクト一覧を表示する",
					Type:        discord.OptionTypeSubCommand,
				},
				{
					Name:        "get",
					Description: "プロジェクトの詳細を取得する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "プロジェクトのID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
					},
				},
				{
					Name:        "delete",
					Description: "プロジェクトを削除する",
					Type:        discord.OptionTypeSubCommand,
					Options: []*discord.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "プロジェクトのID",
							Type:        discord.OptionTypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}

	return b.client.RegisterCommands(commands)
}

// handleInteraction はインタラクションを処理する
func (b *Bot) handleInteraction(i *discord.Interaction) {
	if i.Type != discord.InteractionTypeApplicationCommand {
		return
	}

	if i.Data == nil {
		return
	}

	b.logger.Info("Received command", "name", i.Data.Name)

	switch i.Data.Name {
	case "todo":
		b.handleTodo(i)
	case "task":
		b.handleTask(i)
	case "project":
		b.handleProject(i)
	}
}

// ─── TODO handlers ──────────────────────────────────────

func (b *Bot) handleTodo(i *discord.Interaction) {
	if i.Data == nil || len(i.Data.Options) == 0 {
		return
	}

	subCmd := i.Data.Options[0]
	switch subCmd.Name {
	case "create":
		b.handleTodoCreate(i, subCmd.Options)
	case "list":
		b.handleTodoList(i)
	case "get":
		b.handleTodoGet(i, subCmd.Options)
	case "complete":
		b.handleTodoComplete(i, subCmd.Options)
	case "update":
		b.handleTodoUpdate(i, subCmd.Options)
	case "delete":
		b.handleTodoDelete(i, subCmd.Options)
	}
}

func (b *Bot) handleTodoCreate(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	title := getStringOption(opts, "title")
	description := getStringOption(opts, "description")

	todo, err := b.apiClient.CreateTodo(&api.CreateTodoRequest{
		Title:       title,
		Description: description,
	})
	if err != nil {
		b.respondError(i, "TODO作成に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       "✨ TODO作成完了",
		Description: todo.Title,
		Color:       0x00FF00,
		Fields: []*discord.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", todo.ID), Inline: true},
			{Name: "説明", Value: defaultStr(todo.Description, "なし"), Inline: false},
			{Name: "ステータス", Value: "未完了", Inline: true},
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleTodoList(i *discord.Interaction) {
	todos, err := b.apiClient.ListTodos()
	if err != nil {
		b.respondError(i, "TODO一覧の取得に失敗しました", err)
		return
	}

	if len(todos) == 0 {
		b.respondMessage(i, "📋 TODOはありません")
		return
	}

	var sb strings.Builder
	for idx, todo := range todos {
		status := "⬜"
		if todo.Completed {
			status = "✅"
		}
		sb.WriteString(fmt.Sprintf("%s **%d.** %s\n", status, idx+1, todo.Title))
		sb.WriteString(fmt.Sprintf("   ID: `%s`\n", todo.ID))
	}

	embed := &discord.MessageEmbed{
		Title:       fmt.Sprintf("📋 TODO一覧 (%d件)", len(todos)),
		Description: sb.String(),
		Color:       0x3498DB,
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleTodoGet(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	id := getStringOption(opts, "id")

	todo, err := b.apiClient.GetTodo(id)
	if err != nil {
		b.respondError(i, "TODO取得に失敗しました", err)
		return
	}

	status := "⬜ 未完了"
	if todo.Completed {
		status = "✅ 完了"
	}

	embed := &discord.MessageEmbed{
		Title:       todo.Title,
		Description: defaultStr(todo.Description, "説明なし"),
		Color:       0x3498DB,
		Fields: []*discord.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", todo.ID), Inline: true},
			{Name: "ステータス", Value: status, Inline: true},
			{Name: "作成日", Value: formatTime(todo.CreatedAt), Inline: true},
			{Name: "更新日", Value: formatTime(todo.UpdatedAt), Inline: true},
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleTodoComplete(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	id := getStringOption(opts, "id")
	completed := true

	todo, err := b.apiClient.UpdateTodo(id, &api.UpdateTodoRequest{
		Completed: &completed,
	})
	if err != nil {
		b.respondError(i, "TODO完了処理に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       "✅ TODO完了",
		Description: todo.Title,
		Color:       0x00FF00,
		Fields: []*discord.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", todo.ID), Inline: true},
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleTodoUpdate(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	id := getStringOption(opts, "id")

	req := &api.UpdateTodoRequest{}
	if title := getStringOption(opts, "title"); title != "" {
		req.Title = &title
	}
	if desc := getStringOption(opts, "description"); desc != "" {
		req.Description = &desc
	}

	todo, err := b.apiClient.UpdateTodo(id, req)
	if err != nil {
		b.respondError(i, "TODO更新に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       "📝 TODO更新完了",
		Description: todo.Title,
		Color:       0xF39C12,
		Fields: []*discord.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", todo.ID), Inline: true},
			{Name: "説明", Value: defaultStr(todo.Description, "なし"), Inline: false},
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleTodoDelete(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	id := getStringOption(opts, "id")

	if err := b.apiClient.DeleteTodo(id); err != nil {
		b.respondError(i, "TODO削除に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       "🗑️ TODO削除完了",
		Description: fmt.Sprintf("ID: `%s` のTODOを削除しました", id),
		Color:       0xE74C3C,
	}

	b.respondEmbed(i, embed)
}

// ─── Task handlers ──────────────────────────────────────

func (b *Bot) handleTask(i *discord.Interaction) {
	if i.Data == nil || len(i.Data.Options) == 0 {
		return
	}

	subCmd := i.Data.Options[0]
	switch subCmd.Name {
	case "create":
		b.handleTaskCreate(i, subCmd.Options)
	case "list":
		b.handleTaskList(i, subCmd.Options)
	case "get":
		b.handleTaskGet(i, subCmd.Options)
	case "update":
		b.handleTaskUpdate(i, subCmd.Options)
	case "delete":
		b.handleTaskDelete(i, subCmd.Options)
	}
}

func (b *Bot) handleTaskCreate(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	projectID := getStringOption(opts, "project_id")
	title := getStringOption(opts, "title")
	description := getStringOption(opts, "description")
	status := int(getIntOption(opts, "status", 0))
	priority := int(getIntOption(opts, "priority", 1))

	task, err := b.apiClient.CreateTask(&api.CreateTaskRequest{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
	})
	if err != nil {
		b.respondError(i, "タスク作成に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       "✨ タスク作成完了",
		Description: task.Title,
		Color:       0x00FF00,
		Fields: []*discord.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", task.ID), Inline: true},
			{Name: "プロジェクトID", Value: fmt.Sprintf("`%s`", task.ProjectID), Inline: true},
			{Name: "ステータス", Value: fmt.Sprintf("%s %s", api.StatusEmoji(task.Status), api.StatusText(task.Status)), Inline: true},
			{Name: "優先度", Value: fmt.Sprintf("%s %s", api.PriorityEmoji(task.Priority), api.PriorityText(task.Priority)), Inline: true},
			{Name: "説明", Value: defaultStr(task.Description, "なし"), Inline: false},
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleTaskList(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	projectID := getStringOption(opts, "project_id")

	tasks, err := b.apiClient.ListTasks(projectID)
	if err != nil {
		b.respondError(i, "タスク一覧の取得に失敗しました", err)
		return
	}

	if len(tasks) == 0 {
		b.respondMessage(i, "📋 このプロジェクトにタスクはありません")
		return
	}

	var sb strings.Builder
	for idx, task := range tasks {
		sb.WriteString(fmt.Sprintf("%s %s **%d.** %s\n",
			api.StatusEmoji(task.Status),
			api.PriorityEmoji(task.Priority),
			idx+1,
			task.Title,
		))
		sb.WriteString(fmt.Sprintf("   ID: `%s`\n", task.ID))
	}

	embed := &discord.MessageEmbed{
		Title:       fmt.Sprintf("📋 タスク一覧 (%d件)", len(tasks)),
		Description: sb.String(),
		Color:       0x3498DB,
		Footer: &discord.MessageEmbedFooter{
			Text: fmt.Sprintf("プロジェクトID: %s", projectID),
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleTaskGet(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	id := getStringOption(opts, "id")

	task, err := b.apiClient.GetTask(id)
	if err != nil {
		b.respondError(i, "タスク取得に失敗しました", err)
		return
	}

	endDate := "未設定"
	if task.EndDate != nil {
		endDate = formatTime(*task.EndDate)
	}

	embed := &discord.MessageEmbed{
		Title:       task.Title,
		Description: defaultStr(task.Description, "説明なし"),
		Color:       0x3498DB,
		Fields: []*discord.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", task.ID), Inline: true},
			{Name: "プロジェクトID", Value: fmt.Sprintf("`%s`", task.ProjectID), Inline: true},
			{Name: "ステータス", Value: fmt.Sprintf("%s %s", api.StatusEmoji(task.Status), api.StatusText(task.Status)), Inline: true},
			{Name: "優先度", Value: fmt.Sprintf("%s %s", api.PriorityEmoji(task.Priority), api.PriorityText(task.Priority)), Inline: true},
			{Name: "期限", Value: endDate, Inline: true},
			{Name: "作成日", Value: formatTime(task.CreatedAt), Inline: true},
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleTaskUpdate(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	id := getStringOption(opts, "id")

	// まず既存のタスクを取得
	existing, err := b.apiClient.GetTask(id)
	if err != nil {
		b.respondError(i, "タスク取得に失敗しました", err)
		return
	}

	req := &api.UpdateTaskRequest{
		Title:       existing.Title,
		Description: existing.Description,
		Status:      existing.Status,
		Priority:    existing.Priority,
		EndDate:     existing.EndDate,
	}

	if title := getStringOption(opts, "title"); title != "" {
		req.Title = title
	}
	if desc := getStringOption(opts, "description"); desc != "" {
		req.Description = desc
	}
	if hasOption(opts, "status") {
		req.Status = int(getIntOption(opts, "status", int64(existing.Status)))
	}
	if hasOption(opts, "priority") {
		req.Priority = int(getIntOption(opts, "priority", int64(existing.Priority)))
	}

	task, err := b.apiClient.UpdateTask(id, req)
	if err != nil {
		b.respondError(i, "タスク更新に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       "📝 タスク更新完了",
		Description: task.Title,
		Color:       0xF39C12,
		Fields: []*discord.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", task.ID), Inline: true},
			{Name: "ステータス", Value: fmt.Sprintf("%s %s", api.StatusEmoji(task.Status), api.StatusText(task.Status)), Inline: true},
			{Name: "優先度", Value: fmt.Sprintf("%s %s", api.PriorityEmoji(task.Priority), api.PriorityText(task.Priority)), Inline: true},
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleTaskDelete(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	id := getStringOption(opts, "id")

	if err := b.apiClient.DeleteTask(id); err != nil {
		b.respondError(i, "タスク削除に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       "🗑️ タスク削除完了",
		Description: fmt.Sprintf("ID: `%s` のタスクを削除しました", id),
		Color:       0xE74C3C,
	}

	b.respondEmbed(i, embed)
}

// ─── Project handlers ───────────────────────────────────

func (b *Bot) handleProject(i *discord.Interaction) {
	if i.Data == nil || len(i.Data.Options) == 0 {
		return
	}

	subCmd := i.Data.Options[0]
	switch subCmd.Name {
	case "create":
		b.handleProjectCreate(i, subCmd.Options)
	case "list":
		b.handleProjectList(i)
	case "get":
		b.handleProjectGet(i, subCmd.Options)
	case "delete":
		b.handleProjectDelete(i, subCmd.Options)
	}
}

func (b *Bot) handleProjectCreate(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	title := getStringOption(opts, "title")
	description := getStringOption(opts, "description")

	if b.defaultUserID == "" {
		b.respondError(i, "DEFAULT_USER_ID が設定されていません", fmt.Errorf("環境変数 DEFAULT_USER_ID を設定してください"))
		return
	}

	project, err := b.apiClient.CreateProject(&api.CreateProjectRequest{
		UserID:      b.defaultUserID,
		Title:       title,
		Description: description,
	})
	if err != nil {
		b.respondError(i, "プロジェクト作成に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       "✨ プロジェクト作成完了",
		Description: project.Title,
		Color:       0x00FF00,
		Fields: []*discord.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", project.ID), Inline: true},
			{Name: "説明", Value: defaultStr(project.Description, "なし"), Inline: false},
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleProjectList(i *discord.Interaction) {
	if b.defaultUserID == "" {
		b.respondError(i, "DEFAULT_USER_ID が設定されていません", fmt.Errorf("環境変数 DEFAULT_USER_ID を設定してください"))
		return
	}

	projects, err := b.apiClient.ListProjects(b.defaultUserID)
	if err != nil {
		b.respondError(i, "プロジェクト一覧の取得に失敗しました", err)
		return
	}

	if len(projects) == 0 {
		b.respondMessage(i, "📁 プロジェクトはありません")
		return
	}

	var sb strings.Builder
	for idx, project := range projects {
		sb.WriteString(fmt.Sprintf("📁 **%d.** %s\n", idx+1, project.Title))
		sb.WriteString(fmt.Sprintf("   ID: `%s`\n", project.ID))
		if project.Description != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", project.Description))
		}
	}

	embed := &discord.MessageEmbed{
		Title:       fmt.Sprintf("📁 プロジェクト一覧 (%d件)", len(projects)),
		Description: sb.String(),
		Color:       0x9B59B6,
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleProjectGet(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	id := getStringOption(opts, "id")

	project, err := b.apiClient.GetProject(id)
	if err != nil {
		b.respondError(i, "プロジェクト取得に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       project.Title,
		Description: defaultStr(project.Description, "説明なし"),
		Color:       0x9B59B6,
		Fields: []*discord.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", project.ID), Inline: true},
			{Name: "作成日", Value: formatTime(project.CreatedAt), Inline: true},
			{Name: "更新日", Value: formatTime(project.UpdatedAt), Inline: true},
		},
	}

	b.respondEmbed(i, embed)
}

func (b *Bot) handleProjectDelete(i *discord.Interaction, opts []*discord.InteractionDataOption) {
	id := getStringOption(opts, "id")

	if err := b.apiClient.DeleteProject(id); err != nil {
		b.respondError(i, "プロジェクト削除に失敗しました", err)
		return
	}

	embed := &discord.MessageEmbed{
		Title:       "🗑️ プロジェクト削除完了",
		Description: fmt.Sprintf("ID: `%s` のプロジェクトを削除しました", id),
		Color:       0xE74C3C,
	}

	b.respondEmbed(i, embed)
}

// ─── Response helpers ───────────────────────────────────

func (b *Bot) respondEmbed(i *discord.Interaction, embed *discord.MessageEmbed) {
	err := b.client.RespondToInteraction(i.ID, i.Token, &discord.InteractionResponse{
		Type: discord.InteractionCallbackChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Embeds: []*discord.MessageEmbed{embed},
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond", "error", err)
	}
}

func (b *Bot) respondMessage(i *discord.Interaction, content string) {
	err := b.client.RespondToInteraction(i.ID, i.Token, &discord.InteractionResponse{
		Type: discord.InteractionCallbackChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond", "error", err)
	}
}

func (b *Bot) respondError(i *discord.Interaction, msg string, err error) {
	b.logger.Error(msg, "error", err)
	embed := &discord.MessageEmbed{
		Title:       "❌ エラー",
		Description: fmt.Sprintf("%s\n```%s```", msg, err.Error()),
		Color:       0xE74C3C,
	}
	b.respondEmbed(i, embed)
}

// ─── Option helpers ─────────────────────────────────────

func getStringOption(opts []*discord.InteractionDataOption, name string) string {
	for _, opt := range opts {
		if opt.Name == name {
			return opt.StringValue()
		}
	}
	return ""
}

func getIntOption(opts []*discord.InteractionDataOption, name string, defaultVal int64) int64 {
	for _, opt := range opts {
		if opt.Name == name {
			return opt.IntValue()
		}
	}
	return defaultVal
}

func hasOption(opts []*discord.InteractionDataOption, name string) bool {
	for _, opt := range opts {
		if opt.Name == name {
			return true
		}
	}
	return false
}

func defaultStr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func formatTime(t string) string {
	if len(t) >= 10 {
		return t[:10]
	}
	return t
}

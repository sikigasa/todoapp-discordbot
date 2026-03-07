package bot

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sikigasa/todoapp-discordbot/internal/api"
)

// Bot はDiscord Botを表す
type Bot struct {
	session       *discordgo.Session
	apiClient     *api.Client
	defaultUserID string
	logger        *slog.Logger
}

// New は新しいBotインスタンスを作成する
func New(token string, apiClient *api.Client, defaultUserID string, logger *slog.Logger) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	bot := &Bot{
		session:       session,
		apiClient:     apiClient,
		defaultUserID: defaultUserID,
		logger:        logger,
	}

	session.AddHandler(bot.handleInteraction)
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Info("Bot is ready", "user", s.State.User.Username+"#"+s.State.User.Discriminator)
	})

	return bot, nil
}

// Start はBotを起動し、スラッシュコマンドを登録する
func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
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
	return b.session.Close()
}

// registerCommands はスラッシュコマンドを登録する
func (b *Bot) registerCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "todo",
			Description: "TODOを管理する",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "TODOを作成する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "title",
							Description: "TODOのタイトル",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "description",
							Description: "TODOの説明",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "list",
					Description: "全TODOを一覧表示する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "get",
					Description: "TODOの詳細を取得する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "TODOのID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "complete",
					Description: "TODOを完了にする",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "TODOのID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "update",
					Description: "TODOを更新する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "TODOのID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "title",
							Description: "新しいタイトル",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
						{
							Name:        "description",
							Description: "新しい説明",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "delete",
					Description: "TODOを削除する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "TODOのID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "task",
			Description: "タスクを管理する",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "タスクを作成する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "project_id",
							Description: "プロジェクトID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "title",
							Description: "タスクのタイトル",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "description",
							Description: "タスクの説明",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
						{
							Name:        "status",
							Description: "ステータス (0: To Do, 1: In Progress, 2: Done)",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    false,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{Name: "To Do", Value: 0},
								{Name: "In Progress", Value: 1},
								{Name: "Done", Value: 2},
							},
						},
						{
							Name:        "priority",
							Description: "優先度 (0: Low, 1: Medium, 2: High)",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    false,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "project_id",
							Description: "プロジェクトID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "get",
					Description: "タスクの詳細を取得する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "タスクのID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "update",
					Description: "タスクを更新する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "タスクのID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "title",
							Description: "新しいタイトル",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
						{
							Name:        "description",
							Description: "新しい説明",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
						{
							Name:        "status",
							Description: "ステータス (0: To Do, 1: In Progress, 2: Done)",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    false,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{Name: "To Do", Value: 0},
								{Name: "In Progress", Value: 1},
								{Name: "Done", Value: 2},
							},
						},
						{
							Name:        "priority",
							Description: "優先度 (0: Low, 1: Medium, 2: High)",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    false,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "タスクのID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "project",
			Description: "プロジェクトを管理する",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "プロジェクトを作成する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "title",
							Description: "プロジェクトのタイトル",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "description",
							Description: "プロジェクトの説明",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "list",
					Description: "プロジェクト一覧を表示する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "get",
					Description: "プロジェクトの詳細を取得する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "プロジェクトのID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "delete",
					Description: "プロジェクトを削除する",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "id",
							Description: "プロジェクトのID",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
	}

	for _, cmd := range commands {
		_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("failed to register command %s: %w", cmd.Name, err)
		}
		b.logger.Info("Registered command", "name", cmd.Name)
	}

	return nil
}

// handleInteraction はインタラクションを処理する
func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	b.logger.Info("Received command", "name", data.Name)

	switch data.Name {
	case "todo":
		b.handleTodo(s, i)
	case "task":
		b.handleTask(s, i)
	case "project":
		b.handleProject(s, i)
	}
}

// ─── TODO handlers ──────────────────────────────────────

func (b *Bot) handleTodo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		return
	}

	subCmd := options[0]
	switch subCmd.Name {
	case "create":
		b.handleTodoCreate(s, i, subCmd.Options)
	case "list":
		b.handleTodoList(s, i)
	case "get":
		b.handleTodoGet(s, i, subCmd.Options)
	case "complete":
		b.handleTodoComplete(s, i, subCmd.Options)
	case "update":
		b.handleTodoUpdate(s, i, subCmd.Options)
	case "delete":
		b.handleTodoDelete(s, i, subCmd.Options)
	}
}

func (b *Bot) handleTodoCreate(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	title := getStringOption(opts, "title")
	description := getStringOption(opts, "description")

	todo, err := b.apiClient.CreateTodo(&api.CreateTodoRequest{
		Title:       title,
		Description: description,
	})
	if err != nil {
		b.respondError(s, i, "TODO作成に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✨ TODO作成完了",
		Description: todo.Title,
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", todo.ID), Inline: true},
			{Name: "説明", Value: defaultStr(todo.Description, "なし"), Inline: false},
			{Name: "ステータス", Value: "未完了", Inline: true},
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleTodoList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	todos, err := b.apiClient.ListTodos()
	if err != nil {
		b.respondError(s, i, "TODO一覧の取得に失敗しました", err)
		return
	}

	if len(todos) == 0 {
		b.respondMessage(s, i, "📋 TODOはありません")
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

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("📋 TODO一覧 (%d件)", len(todos)),
		Description: sb.String(),
		Color:       0x3498DB,
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleTodoGet(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	id := getStringOption(opts, "id")

	todo, err := b.apiClient.GetTodo(id)
	if err != nil {
		b.respondError(s, i, "TODO取得に失敗しました", err)
		return
	}

	status := "⬜ 未完了"
	if todo.Completed {
		status = "✅ 完了"
	}

	embed := &discordgo.MessageEmbed{
		Title:       todo.Title,
		Description: defaultStr(todo.Description, "説明なし"),
		Color:       0x3498DB,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", todo.ID), Inline: true},
			{Name: "ステータス", Value: status, Inline: true},
			{Name: "作成日", Value: formatTime(todo.CreatedAt), Inline: true},
			{Name: "更新日", Value: formatTime(todo.UpdatedAt), Inline: true},
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleTodoComplete(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	id := getStringOption(opts, "id")
	completed := true

	todo, err := b.apiClient.UpdateTodo(id, &api.UpdateTodoRequest{
		Completed: &completed,
	})
	if err != nil {
		b.respondError(s, i, "TODO完了処理に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✅ TODO完了",
		Description: todo.Title,
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", todo.ID), Inline: true},
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleTodoUpdate(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
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
		b.respondError(s, i, "TODO更新に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📝 TODO更新完了",
		Description: todo.Title,
		Color:       0xF39C12,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", todo.ID), Inline: true},
			{Name: "説明", Value: defaultStr(todo.Description, "なし"), Inline: false},
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleTodoDelete(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	id := getStringOption(opts, "id")

	if err := b.apiClient.DeleteTodo(id); err != nil {
		b.respondError(s, i, "TODO削除に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🗑️ TODO削除完了",
		Description: fmt.Sprintf("ID: `%s` のTODOを削除しました", id),
		Color:       0xE74C3C,
	}

	b.respondEmbed(s, i, embed)
}

// ─── Task handlers ──────────────────────────────────────

func (b *Bot) handleTask(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		return
	}

	subCmd := options[0]
	switch subCmd.Name {
	case "create":
		b.handleTaskCreate(s, i, subCmd.Options)
	case "list":
		b.handleTaskList(s, i, subCmd.Options)
	case "get":
		b.handleTaskGet(s, i, subCmd.Options)
	case "update":
		b.handleTaskUpdate(s, i, subCmd.Options)
	case "delete":
		b.handleTaskDelete(s, i, subCmd.Options)
	}
}

func (b *Bot) handleTaskCreate(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
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
		b.respondError(s, i, "タスク作成に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✨ タスク作成完了",
		Description: task.Title,
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", task.ID), Inline: true},
			{Name: "プロジェクトID", Value: fmt.Sprintf("`%s`", task.ProjectID), Inline: true},
			{Name: "ステータス", Value: fmt.Sprintf("%s %s", api.StatusEmoji(task.Status), api.StatusText(task.Status)), Inline: true},
			{Name: "優先度", Value: fmt.Sprintf("%s %s", api.PriorityEmoji(task.Priority), api.PriorityText(task.Priority)), Inline: true},
			{Name: "説明", Value: defaultStr(task.Description, "なし"), Inline: false},
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleTaskList(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	projectID := getStringOption(opts, "project_id")

	tasks, err := b.apiClient.ListTasks(projectID)
	if err != nil {
		b.respondError(s, i, "タスク一覧の取得に失敗しました", err)
		return
	}

	if len(tasks) == 0 {
		b.respondMessage(s, i, "📋 このプロジェクトにタスクはありません")
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

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("📋 タスク一覧 (%d件)", len(tasks)),
		Description: sb.String(),
		Color:       0x3498DB,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("プロジェクトID: %s", projectID),
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleTaskGet(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	id := getStringOption(opts, "id")

	task, err := b.apiClient.GetTask(id)
	if err != nil {
		b.respondError(s, i, "タスク取得に失敗しました", err)
		return
	}

	endDate := "未設定"
	if task.EndDate != nil {
		endDate = formatTime(*task.EndDate)
	}

	embed := &discordgo.MessageEmbed{
		Title:       task.Title,
		Description: defaultStr(task.Description, "説明なし"),
		Color:       0x3498DB,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", task.ID), Inline: true},
			{Name: "プロジェクトID", Value: fmt.Sprintf("`%s`", task.ProjectID), Inline: true},
			{Name: "ステータス", Value: fmt.Sprintf("%s %s", api.StatusEmoji(task.Status), api.StatusText(task.Status)), Inline: true},
			{Name: "優先度", Value: fmt.Sprintf("%s %s", api.PriorityEmoji(task.Priority), api.PriorityText(task.Priority)), Inline: true},
			{Name: "期限", Value: endDate, Inline: true},
			{Name: "作成日", Value: formatTime(task.CreatedAt), Inline: true},
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleTaskUpdate(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	id := getStringOption(opts, "id")

	// まず既存のタスクを取得
	existing, err := b.apiClient.GetTask(id)
	if err != nil {
		b.respondError(s, i, "タスク取得に失敗しました", err)
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
		b.respondError(s, i, "タスク更新に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📝 タスク更新完了",
		Description: task.Title,
		Color:       0xF39C12,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", task.ID), Inline: true},
			{Name: "ステータス", Value: fmt.Sprintf("%s %s", api.StatusEmoji(task.Status), api.StatusText(task.Status)), Inline: true},
			{Name: "優先度", Value: fmt.Sprintf("%s %s", api.PriorityEmoji(task.Priority), api.PriorityText(task.Priority)), Inline: true},
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleTaskDelete(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	id := getStringOption(opts, "id")

	if err := b.apiClient.DeleteTask(id); err != nil {
		b.respondError(s, i, "タスク削除に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🗑️ タスク削除完了",
		Description: fmt.Sprintf("ID: `%s` のタスクを削除しました", id),
		Color:       0xE74C3C,
	}

	b.respondEmbed(s, i, embed)
}

// ─── Project handlers ───────────────────────────────────

func (b *Bot) handleProject(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		return
	}

	subCmd := options[0]
	switch subCmd.Name {
	case "create":
		b.handleProjectCreate(s, i, subCmd.Options)
	case "list":
		b.handleProjectList(s, i)
	case "get":
		b.handleProjectGet(s, i, subCmd.Options)
	case "delete":
		b.handleProjectDelete(s, i, subCmd.Options)
	}
}

func (b *Bot) handleProjectCreate(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	title := getStringOption(opts, "title")
	description := getStringOption(opts, "description")

	if b.defaultUserID == "" {
		b.respondError(s, i, "DEFAULT_USER_ID が設定されていません", fmt.Errorf("環境変数 DEFAULT_USER_ID を設定してください"))
		return
	}

	project, err := b.apiClient.CreateProject(&api.CreateProjectRequest{
		UserID:      b.defaultUserID,
		Title:       title,
		Description: description,
	})
	if err != nil {
		b.respondError(s, i, "プロジェクト作成に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✨ プロジェクト作成完了",
		Description: project.Title,
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", project.ID), Inline: true},
			{Name: "説明", Value: defaultStr(project.Description, "なし"), Inline: false},
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleProjectList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.defaultUserID == "" {
		b.respondError(s, i, "DEFAULT_USER_ID が設定されていません", fmt.Errorf("環境変数 DEFAULT_USER_ID を設定してください"))
		return
	}

	projects, err := b.apiClient.ListProjects(b.defaultUserID)
	if err != nil {
		b.respondError(s, i, "プロジェクト一覧の取得に失敗しました", err)
		return
	}

	if len(projects) == 0 {
		b.respondMessage(s, i, "📁 プロジェクトはありません")
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

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("📁 プロジェクト一覧 (%d件)", len(projects)),
		Description: sb.String(),
		Color:       0x9B59B6,
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleProjectGet(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	id := getStringOption(opts, "id")

	project, err := b.apiClient.GetProject(id)
	if err != nil {
		b.respondError(s, i, "プロジェクト取得に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       project.Title,
		Description: defaultStr(project.Description, "説明なし"),
		Color:       0x9B59B6,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID", Value: fmt.Sprintf("`%s`", project.ID), Inline: true},
			{Name: "作成日", Value: formatTime(project.CreatedAt), Inline: true},
			{Name: "更新日", Value: formatTime(project.UpdatedAt), Inline: true},
		},
	}

	b.respondEmbed(s, i, embed)
}

func (b *Bot) handleProjectDelete(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	id := getStringOption(opts, "id")

	if err := b.apiClient.DeleteProject(id); err != nil {
		b.respondError(s, i, "プロジェクト削除に失敗しました", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🗑️ プロジェクト削除完了",
		Description: fmt.Sprintf("ID: `%s` のプロジェクトを削除しました", id),
		Color:       0xE74C3C,
	}

	b.respondEmbed(s, i, embed)
}

// ─── Response helpers ───────────────────────────────────

func (b *Bot) respondEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond", "error", err)
	}
}

func (b *Bot) respondMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond", "error", err)
	}
}

func (b *Bot) respondError(s *discordgo.Session, i *discordgo.InteractionCreate, msg string, err error) {
	b.logger.Error(msg, "error", err)
	embed := &discordgo.MessageEmbed{
		Title:       "❌ エラー",
		Description: fmt.Sprintf("%s\n```%s```", msg, err.Error()),
		Color:       0xE74C3C,
	}
	b.respondEmbed(s, i, embed)
}

// ─── Option helpers ─────────────────────────────────────

func getStringOption(opts []*discordgo.ApplicationCommandInteractionDataOption, name string) string {
	for _, opt := range opts {
		if opt.Name == name {
			return opt.StringValue()
		}
	}
	return ""
}

func getIntOption(opts []*discordgo.ApplicationCommandInteractionDataOption, name string, defaultVal int64) int64 {
	for _, opt := range opts {
		if opt.Name == name {
			return opt.IntValue()
		}
	}
	return defaultVal
}

func hasOption(opts []*discordgo.ApplicationCommandInteractionDataOption, name string) bool {
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

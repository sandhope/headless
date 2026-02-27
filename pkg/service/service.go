package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func GetCookie(domain string) (result map[string]string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// 基础运行参数
		chromedp.Flag("headless", true),              // 必需 - 无头模式运行
		chromedp.Flag("disable-gpu", true),           // 推荐 - 禁用GPU加速（无头模式下通常不需要）
		chromedp.Flag("disable-cache", true),         // 可选 - 禁用缓存，确保获取最新数据
		chromedp.Flag("no-sandbox", true),            // 必需 - 禁用沙盒模式（容器环境中必需）
		chromedp.Flag("disable-dev-shm-usage", true), // 必需 - 禁用/dev/shm使用（Docker环境常见问题）

		// 扩展功能禁用
		chromedp.Flag("disable-extensions", true),                     // 推荐 - 禁用扩展程序
		chromedp.Flag("disable-background-timer-throttling", true),    // 可选 - 禁用后台定时器节流
		chromedp.Flag("disable-backgrounding-occluded-windows", true), // 可选 - 禁用被遮挡窗口的后台处理
		chromedp.Flag("disable-renderer-backgrounding", true),         // 可选 - 禁用渲染器后台处理

		// 反检测相关参数
		chromedp.Flag("disable-blink-features", "AutomationControlled"),                                                                       // 必需 - 禁用自动化检测特征
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"), // 必需 - 设置真实用户代理
		chromedp.WindowSize(1920, 1080),                 // 推荐 - 设置合理窗口大小
		chromedp.Flag("no-default-browser-check", true), // 可选 - 跳过默认浏览器检查
		chromedp.NoFirstRun,                             // 推荐 - 跳过首次运行向导
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	var cookies []*network.Cookie

	tasks := chromedp.Tasks{
		// 在导航前，注入 JS 抹除 webdriver 特征
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	_, _, err := runtime.Evaluate(`Object.defineProperty(navigator, 'webdriver', {get: () => undefined})`).Do(ctx)
		// 	return err
		// }),

		chromedp.Navigate(fmt.Sprintf("https://%s", domain)),
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	}

	if err = chromedp.Run(ctx, tasks); err != nil {
		return nil, fmt.Errorf("failed to get cookies for domain %s: %w", domain, err)
	}

	result = make(map[string]string, len(cookies))
	for _, cookie := range cookies {
		result[cookie.Name] = cookie.Value
	}

	return
}

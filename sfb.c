#include "sfb.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int sfb_init(SFB_Context *context, const char *title, int x, int y, int width,
             int height, uint32_t flags) {
  if (SDL_Init(SDL_INIT_VIDEO) < 0) {
    printf("SDL could not initialize! SDL_Error: %s\n", SDL_GetError());
    return -1;
  }

  context->window = SDL_CreateWindow(title, x, y, width, height, flags);
  if (context->window == NULL) {
    printf("Window could not be created! SDL_Error: %s\n", SDL_GetError());
    SDL_Quit();
    return -1;
  }

  context->renderer =
      SDL_CreateRenderer(context->window, -1, SDL_RENDERER_ACCELERATED);
  if (context->renderer == NULL) {
    printf("Accelerated renderer could not be created! Falling back to "
           "software renderer. SDL_Error: %s\n",
           SDL_GetError());
    context->renderer =
        SDL_CreateRenderer(context->window, -1, SDL_RENDERER_SOFTWARE);
    if (context->renderer == NULL) {
      printf("Software renderer could not be created! SDL_Error: %s\n",
             SDL_GetError());
      SDL_DestroyWindow(context->window);
      SDL_Quit();
      return -1;
    }
  }

  context->texture =
      SDL_CreateTexture(context->renderer, SDL_PIXELFORMAT_RGBA8888,
                        SDL_TEXTUREACCESS_STREAMING, width, height);
  if (context->texture == NULL) {
    printf("Texture could not be created! SDL_Error: %s\n", SDL_GetError());
    SDL_DestroyRenderer(context->renderer);
    SDL_DestroyWindow(context->window);
    SDL_Quit();
    return -1;
  }

  context->framebuffer = (Uint32 *)malloc(width * height * sizeof(Uint32));
  if (context->framebuffer == NULL) {
    printf("Framebuffer could not be allocated!\n");
    SDL_DestroyTexture(context->texture);
    SDL_DestroyRenderer(context->renderer);
    SDL_DestroyWindow(context->window);
    SDL_Quit();
    return -1;
  }

  context->width = width;
  context->height = height;
  context->last_time = SDL_GetTicks();
  context->frame_count = 0;

  context->pixel_format = SDL_AllocFormat(SDL_PIXELFORMAT_RGBA8888);

  return 0;
}

void sfb_destroy(SFB_Context *context) {
  if (context->framebuffer) {
    free(context->framebuffer);
    context->framebuffer = NULL;
  }
  if (context->texture) {
    SDL_DestroyTexture(context->texture);
    context->texture = NULL;
  }
  if (context->renderer) {
    SDL_DestroyRenderer(context->renderer);
    context->renderer = NULL;
  }
  if (context->window) {
    SDL_DestroyWindow(context->window);
    context->window = NULL;
  }
  if (context->pixel_format) {
    SDL_FreeFormat(context->pixel_format);
    context->pixel_format = NULL;
  }
  SDL_Quit();
}

void sfb_render(SFB_Context *context) {
  SDL_UpdateTexture(context->texture, NULL, context->framebuffer,
                    context->width * sizeof(Uint32));
  SDL_RenderClear(context->renderer);
  SDL_RenderCopy(context->renderer, context->texture, NULL, NULL);
  SDL_RenderPresent(context->renderer);
}

void sfb_resize(SFB_Context *context, int width, int height) {
  if (width <= 0 || height <= 0) {
    printf("Invalid dimensions for resize: width=%d, height=%d\n", width,
           height);
    return;
  }

  float aspect_ratio = (float)context->width / context->height;
  int new_width, new_height;

  if (width / (float)height > aspect_ratio) {
    new_height = height;
    new_width = (int)(height * aspect_ratio);
  } else {
    new_width = width;
    new_height = (int)(width / aspect_ratio);
  }

  SDL_SetWindowSize(context->window, new_width, new_height);

  Uint32 *new_framebuffer =
      (Uint32 *)realloc(context->framebuffer, width * height * sizeof(Uint32));
  if (new_framebuffer == NULL) {
    printf("Framebuffer could not be reallocated!\n");
    return;
  }
  context->framebuffer = new_framebuffer;

  SDL_DestroyTexture(context->texture);
  context->texture =
      SDL_CreateTexture(context->renderer, SDL_PIXELFORMAT_RGBA8888,
                        SDL_TEXTUREACCESS_STREAMING, new_width, new_height);
  if (context->texture == NULL) {
    printf("Texture could not be recreated! SDL_Error: %s\n", SDL_GetError());
    return;
  }

  context->width = new_width;
  context->height = new_height;
}

int sfb_poll_event(uint32_t *eventType, uint32_t *keySym) {
  SDL_Event e;
  if (SDL_PollEvent(&e) != 0) {
    *eventType = e.type;
    if (e.type == SDL_KEYDOWN || e.type == SDL_KEYUP) {
      *keySym = e.key.keysym.sym;
    }
    return 1;
  }
  return 0;
}

void fill_framebuffer(Uint32 *framebuffer, Pixel color, int width, int height) {
  SDL_PixelFormat *format = SDL_AllocFormat(SDL_PIXELFORMAT_RGBA8888);
  Uint32 pixel_value = SDL_MapRGBA(format, color.r, color.g, color.b, color.a);
  for (int y = 0; y < height; y++) {
    for (int x = 0; x < width; x++) {
      framebuffer[y * width + x] = pixel_value;
    }
  }
  SDL_FreeFormat(format);
}

void sfb_set_pixel(SFB_Context *context, int x, int y, Pixel color) {
  if (x >= 0 && x < context->width && y >= 0 && y < context->height) {
    context->framebuffer[y * context->width + x] =
        SDL_MapRGBA(context->pixel_format, color.r, color.g, color.b, color.a);
  }
}

void sfb_set_window_title(SFB_Context *context, const char *title) {
  if (context && context->window) {
    if (context->original_title[0] == '\0') {
      strncpy(context->original_title, title,
              sizeof(context->original_title) - 1);
      context->original_title[sizeof(context->original_title) - 1] = '\0';
    }
    SDL_SetWindowTitle(context->window, title);
  }
}

void sfb_enable_fps(SFB_Context *context) {
  context->frame_count++;
  uint32_t current_time = SDL_GetTicks();
  uint32_t elapsed_time = current_time - context->last_time;

  if (elapsed_time > 1000) {
    float fps = context->frame_count / (elapsed_time / 1000.0f);
    char title[256];
    size_t fps_string_length = snprintf(NULL, 0, " - FPS: %.2f", fps);
    size_t max_title_length = sizeof(title) - fps_string_length - 1;

    char truncated_title[256];
    strncpy(truncated_title, context->original_title, max_title_length);
    truncated_title[max_title_length] = '\0';

    snprintf(title, sizeof(title), "%.*s - FPS: %.2f", (int)max_title_length,
             truncated_title, fps);
    sfb_set_window_title(context, title);

    context->last_time = current_time;
    context->frame_count = 0;
  }
}

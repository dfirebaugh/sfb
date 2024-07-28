#ifndef SFB_H
#define SFB_H

#include <SDL2/SDL.h>
#include <stdint.h>

typedef struct {
  uint8_t r;
  uint8_t g;
  uint8_t b;
  uint8_t a;
} Pixel;

typedef struct {
  SDL_Window *window;
  SDL_Renderer *renderer;
  SDL_Texture *texture;
  SDL_PixelFormat *pixel_format;
  Uint32 *framebuffer;
  int width;
  int height;
  uint32_t frame_count;
  uint32_t last_time;
  char original_title[256];
} SFB_Context;

int sfb_init(SFB_Context *context, const char *title, int x, int y, int width,
             int height, uint32_t flags);
void sfb_destroy(SFB_Context *context);
void sfb_render(SFB_Context *context);
void sfb_resize(SFB_Context *context, int width, int height);
int sfb_poll_event(uint32_t *eventType, uint32_t *keySym);
void fill_framebuffer(Uint32 *framebuffer, Pixel color, int width, int height);
void sfb_set_pixel(SFB_Context *context, int x, int y, Pixel color);
void sfb_set_window_title(SFB_Context *context, const char *title);
void sfb_enable_fps(SFB_Context *context);

#endif // SFB_H

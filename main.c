#if 0
#include <SDL2/SDL.h>
#include <stdbool.h>
#include <stdint.h>

#include "sfb.h"
void draw_triangle(SFB_Context *context, Pixel color, int x1, int y1, int x2,
                   int y2, int x3, int y3) {
  if (y1 > y2) {
    int temp;
    temp = y1;
    y1 = y2;
    y2 = temp;
    temp = x1;
    x1 = x2;
    x2 = temp;
  }
  if (y2 > y3) {
    int temp;
    temp = y2;
    y2 = y3;
    y3 = temp;
    temp = x2;
    x2 = x3;
    x3 = temp;
  }
  if (y1 > y2) {
    int temp;
    temp = y1;
    y1 = y2;
    y2 = temp;
    temp = x1;
    x1 = x2;
    x2 = temp;
  }

  int total_height = y3 - y1;
  for (int i = 0; i < total_height; i++) {
    bool second_half = i > y2 - y1 || y2 == y1;
    int segment_height = second_half ? y3 - y2 : y2 - y1;
    float alpha = (float)i / total_height;
    float beta = (float)(i - (second_half ? y2 - y1 : 0)) / segment_height;
    int A_x = x1 + (x3 - x1) * alpha;
    int B_x = second_half ? x2 + (x3 - x2) * beta : x1 + (x2 - x1) * beta;
    int A_y = y1 + i;
    int B_y = second_half ? y2 + (y3 - y2) * beta : y1 + i;

    if (A_x > B_x) {
      int temp = A_x;
      A_x = B_x;
      B_x = temp;
    }

    for (int j = A_x; j <= B_x; j++) {
      sfb_set_pixel(context, j, A_y, color);
    }
  }
}

int main(int argc, char *args[]) {
  SFB_Context context;

  if (sfb_init(&context, "SDL2 Triangle", SDL_WINDOWPOS_UNDEFINED,
               SDL_WINDOWPOS_UNDEFINED, 800, 600, SDL_WINDOW_SHOWN) != 0) {
    return 1;
  }

  Pixel color = {.r = 255, .g = 255, .b = 0, .a = 255}; // Yellow color

  int quit = 0;
  uint32_t eventType;
  uint32_t keySym;

  while (!quit) {
    while (sfb_poll_event(&eventType, &keySym)) {
      if (eventType == SDL_QUIT ||
          (eventType == SDL_KEYDOWN && keySym == SDLK_ESCAPE)) {
        quit = 1;
      }
    }

    fill_framebuffer(context.framebuffer,
                     (Pixel){.r = 0, .g = 0, .b = 0, .a = 255},
                     context.width, context.height);
                     black
    pixel draw_triangle(&context, color, 400, 150, 200, 450, 600, 450);

    printf("First few pixels before update:\n");
    for (int i = 0; i < 10; i++) {
      printf("Pixel %d: R=%d, G=%d, B=%d, A=%d\n", i,
      context.framebuffer[i].r,
             context.framebuffer[i].g, context.framebuffer[i].b,
             context.framebuffer[i].a);
    }

    uint32_t format;
    int access, w, h;
    SDL_QueryTexture(context.texture, &format, &access, &w, &h);
    printf("Texture Pixel Format: %s\n", SDL_GetPixelFormatName(format));

    sfb_render(&context);
  }

  sfb_destroy(&context);
  return 0;
}
#endif
